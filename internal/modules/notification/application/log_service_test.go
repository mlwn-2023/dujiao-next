package application

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/dujiao-next/internal/constants"
	"github.com/dujiao-next/internal/modules/notification/contract"
	"github.com/dujiao-next/internal/modules/notification/domain"
	settingsmessaging "github.com/dujiao-next/internal/modules/settings/schema/messaging"
	settingsstorefront "github.com/dujiao-next/internal/modules/settings/schema/storefront"
	"github.com/dujiao-next/internal/queue"
)

type notificationSettingsStub struct {
	notification settingsmessaging.NotificationCenterSetting
	dashboard    settingsstorefront.DashboardSetting
}

func (s notificationSettingsStub) GetNotificationCenterSetting() (settingsmessaging.NotificationCenterSetting, error) {
	return s.notification, nil
}

func (s notificationSettingsStub) GetDashboardSetting() (settingsstorefront.DashboardSetting, error) {
	return s.dashboard, nil
}

type notificationEmailStub struct{}

func (notificationEmailStub) SendCustomEmail(recipient, _, _ string) error {
	if strings.Contains(recipient, "failure") {
		return errors.New("simulated email failure")
	}
	return nil
}

type notificationWXPushStub struct {
	inputs []contract.WXPushSendInput
}

func (s *notificationWXPushStub) Send(_ context.Context, input contract.WXPushSendInput) error {
	s.inputs = append(s.inputs, input)
	return nil
}

type notificationLogRepositoryStub struct {
	items []domain.NotificationLog
}

func (r *notificationLogRepositoryStub) Create(item *domain.NotificationLog) error {
	if item != nil {
		copy := *item
		copy.ID = uint(len(r.items) + 1)
		r.items = append(r.items, copy)
	}
	return nil
}

func (r *notificationLogRepositoryStub) ListAdmin(filter contract.LogListFilter) ([]domain.NotificationLog, int64, error) {
	result := make([]domain.NotificationLog, 0, len(r.items))
	for _, item := range r.items {
		if filter.EventType != "" && item.EventType != filter.EventType {
			continue
		}
		if filter.IsTest != nil && item.IsTest != *filter.IsTest {
			continue
		}
		result = append(result, item)
	}
	return result, int64(len(result)), nil
}

func setupLogService(t *testing.T) (*Service, *LogService) {
	t.Helper()
	logService := NewLogService(&notificationLogRepositoryStub{})
	setting := settingsmessaging.NotificationCenterSetting{
		DefaultLocale: constants.LocaleEnUS,
		Scenes: settingsmessaging.NotificationSceneSetting{
			OrderPaidSuccess: true,
			ExceptionAlert:   true,
		},
		Channels: settingsmessaging.NotificationChannelsSetting{
			Email: settingsmessaging.NotificationChannelSetting{
				Enabled:    true,
				Recipients: []string{"success@example.com", "failure@example.com"},
			},
		},
		Templates: settingsmessaging.NotificationTemplatesSetting{
			OrderPaidSuccess: settingsmessaging.NotificationSceneTemplate{
				ENUS: settingsmessaging.NotificationLocalizedTemplate{
					Title: "Order {{order_no}}",
					Body:  "Customer {{customer_email}}",
				},
			},
			ExceptionAlert: settingsmessaging.NotificationSceneTemplate{
				ENUS: settingsmessaging.NotificationLocalizedTemplate{
					Title: "Alert {{message}}",
					Body:  "{{message}}",
				},
			},
		},
	}
	service := NewService(notificationSettingsStub{notification: setting}, notificationEmailStub{}, nil, nil, logService, nil, nil)
	return service, logService
}

func TestServiceSendTestUsesWXPushSettingsAndGroupOverride(t *testing.T) {
	setting := settingsmessaging.NotificationCenterDefaultSetting()
	setting.DefaultLocale = constants.LocaleZhCN
	setting.Channels.WXPush = settingsmessaging.NotificationWXPushChannelSetting{
		Enabled:  true,
		BaseURL:  "https://push.example.com",
		APIToken: "secret-token",
		Groups:   []string{"默认分组"},
	}

	gateway := &notificationWXPushStub{}
	logService := NewLogService(&notificationLogRepositoryStub{})
	service := NewService(notificationSettingsStub{notification: setting}, nil, nil, nil, logService, nil, gateway)
	err := service.SendTest(context.Background(), contract.TestSendInput{
		Channel: "wxpush",
		Target:  "服务器告警|管理员",
		Scene:   constants.NotificationEventOrderPaidSuccess,
		Locale:  constants.LocaleZhCN,
	})
	if err != nil {
		t.Fatalf("SendTest() error = %v", err)
	}
	if len(gateway.inputs) != 1 {
		t.Fatalf("expected one wxpush request, got %d", len(gateway.inputs))
	}
	input := gateway.inputs[0]
	if input.BaseURL != "https://push.example.com" || input.APIToken != "secret-token" {
		t.Fatalf("unexpected connection settings: %#v", input)
	}
	if len(input.Groups) != 2 || input.Groups[0] != "服务器告警" || input.Groups[1] != "管理员" {
		t.Fatalf("unexpected groups: %#v", input.Groups)
	}
	if input.Title == "" || input.Body == "" {
		t.Fatalf("expected rendered message, got %#v", input)
	}

	isTest := true
	logs, total, err := logService.ListForAdmin(contract.LogListFilter{Page: 1, PageSize: 10, IsTest: &isTest})
	if err != nil || total != 1 || len(logs) != 1 {
		t.Fatalf("unexpected logs: total=%d len=%d err=%v", total, len(logs), err)
	}
	if logs[0].Channel != "wxpush" || logs[0].Recipient != "服务器告警 | 管理员" {
		t.Fatalf("unexpected log: %#v", logs[0])
	}
}

func TestServiceSendTestRecordsSuccessLog(t *testing.T) {
	service, logService := setupLogService(t)
	if err := service.SendTest(context.Background(), contract.TestSendInput{
		Channel: "email",
		Target:  "success@example.com",
		Scene:   constants.NotificationEventOrderPaidSuccess,
		Locale:  constants.LocaleEnUS,
	}); err != nil {
		t.Fatalf("SendTest failed: %v", err)
	}

	isTest := true
	items, total, err := logService.ListForAdmin(contract.LogListFilter{
		Page: 1, PageSize: 10, EventType: constants.NotificationEventOrderPaidSuccess, IsTest: &isTest,
	})
	if err != nil {
		t.Fatalf("list notification logs failed: %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("expected 1 test log, total=%d len=%d", total, len(items))
	}
	if items[0].Status != notificationLogStatusSuccess || !items[0].IsTest {
		t.Fatalf("unexpected log: %#v", items[0])
	}
	if !strings.Contains(items[0].Title, "DJ202603230001") {
		t.Fatalf("title should include sample order no, got %s", items[0].Title)
	}
}

func TestServiceDispatchSingleEventRecordsPerRecipientResult(t *testing.T) {
	service, logService := setupLogService(t)
	setting, err := service.settingService.GetNotificationCenterSetting()
	if err != nil {
		t.Fatalf("get notification center setting failed: %v", err)
	}
	dispatchErr := service.dispatchSingleEvent(context.Background(), setting, queue.NotificationDispatchPayload{
		EventType: constants.NotificationEventOrderPaidSuccess,
		BizType:   constants.NotificationBizTypeOrder,
		BizID:     88,
		Locale:    constants.LocaleEnUS,
		Force:     true,
		Data: map[string]interface{}{
			"order_no":       "DJ-LOG-88",
			"customer_email": "member@example.com",
		},
	})
	if !errors.Is(dispatchErr, contract.ErrSendFailed) {
		t.Fatalf("expected send failure, got %v", dispatchErr)
	}

	items, total, err := logService.ListForAdmin(contract.LogListFilter{Page: 1, PageSize: 10, EventType: constants.NotificationEventOrderPaidSuccess})
	if err != nil {
		t.Fatalf("list notification logs failed: %v", err)
	}
	if total != 2 || len(items) != 2 {
		t.Fatalf("expected 2 recipient logs, total=%d len=%d", total, len(items))
	}
	statuses := map[string]string{}
	for _, item := range items {
		statuses[item.Recipient] = item.Status
	}
	if statuses["success@example.com"] != notificationLogStatusSuccess {
		t.Fatalf("success recipient status mismatch: %v", statuses)
	}
	if statuses["failure@example.com"] != notificationLogStatusFailed {
		t.Fatalf("failure recipient status mismatch: %v", statuses)
	}
}
