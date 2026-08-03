package settingsmessaging

import "testing"

func TestApplyNotificationCenterSettingPatchKeepsMaskedWXPushToken(t *testing.T) {
	current := NotificationCenterDefaultSetting()
	current.Channels.WXPush.APIToken = "existing-token"

	enabled := true
	baseURL := " https://push.example.com/ "
	emptyToken := ""
	groups := []string{"服务器告警", "服务器告警", " 管理员 "}
	next, err := ApplyNotificationCenterSettingPatch(current, NotificationCenterSettingPatch{
		Channels: &NotificationChannelsPatch{
			WXPush: &NotificationWXPushChannelPatch{
				Enabled:  &enabled,
				BaseURL:  &baseURL,
				APIToken: &emptyToken,
				Groups:   &groups,
			},
		},
	})
	if err != nil {
		t.Fatalf("ApplyNotificationCenterSettingPatch() error = %v", err)
	}
	if next.Channels.WXPush.APIToken != "existing-token" {
		t.Fatalf("token should be preserved, got %q", next.Channels.WXPush.APIToken)
	}
	if next.Channels.WXPush.BaseURL != "https://push.example.com" {
		t.Fatalf("unexpected base url: %q", next.Channels.WXPush.BaseURL)
	}
	if len(next.Channels.WXPush.Groups) != 2 || next.Channels.WXPush.Groups[1] != "管理员" {
		t.Fatalf("unexpected groups: %#v", next.Channels.WXPush.Groups)
	}

	masked := MaskNotificationCenterSettingForAdmin(next)
	channels, ok := masked["channels"].(map[string]interface{})
	if !ok {
		t.Fatalf("missing channels: %#v", masked)
	}
	wxpush, ok := channels["wxpush"].(map[string]interface{})
	if !ok {
		t.Fatalf("missing wxpush channel: %#v", channels)
	}
	if wxpush["api_token"] != "" || wxpush["has_api_token"] != true {
		t.Fatalf("wxpush token was not masked: %#v", wxpush)
	}
}

func TestValidateNotificationCenterSettingRejectsInvalidWXPushURL(t *testing.T) {
	setting := NotificationCenterDefaultSetting()
	setting.Channels.WXPush = NotificationWXPushChannelSetting{
		Enabled:  true,
		BaseURL:  "file:///tmp/wxpush",
		APIToken: "token",
	}
	if err := ValidateNotificationCenterSetting(setting); err == nil {
		t.Fatal("ValidateNotificationCenterSetting() expected error")
	}
}
