package wxpush

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dujiao-next/internal/modules/notification/contract"
)

func TestClientSendSignsAndPostsWXPushRequest(t *testing.T) {
	const token = "wxpush-test-token"
	fixedNow := time.Date(2026, 8, 3, 7, 30, 0, 123_000_000, time.UTC)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/wxsend" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer "+token {
			t.Fatalf("unexpected authorization header: %q", got)
		}
		var body sendRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if body.Timestamp != "1785742200123" {
			t.Fatalf("unexpected timestamp: %s", body.Timestamp)
		}
		mac := hmac.New(sha256.New, []byte(token))
		_, _ = mac.Write([]byte(body.Timestamp + "\n" + token))
		wantSign := base64.StdEncoding.EncodeToString(mac.Sum(nil))
		if body.Sign != wantSign {
			t.Fatalf("unexpected signature: %q", body.Sign)
		}
		if len(body.Groups) != 2 || body.Groups[0] != "服务器告警" || body.Groups[1] != "管理员" {
			t.Fatalf("unexpected groups: %#v", body.Groups)
		}
		if body.Title != "订单支付成功" || body.Content != "订单 DJ-1001 已支付" {
			t.Fatalf("unexpected message: title=%q content=%q", body.Title, body.Content)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"msg":"Successfully sent messages to 2 user(s)."}`))
	}))
	defer server.Close()

	client := New(server.Client())
	client.now = func() time.Time { return fixedNow }
	err := client.Send(t.Context(), contract.WXPushSendInput{
		BaseURL:  server.URL,
		APIToken: token,
		Groups:   []string{"服务器告警", "管理员"},
		Title:    "订单支付成功",
		Body:     "订单 DJ-1001 已支付",
	})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
}

func TestClientSendAcceptsFullWXSendEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/prefix/wxsend" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(server.Client())
	if err := client.Send(t.Context(), contract.WXPushSendInput{
		BaseURL:  server.URL + "/prefix/wxsend",
		APIToken: "token",
		Title:    "title",
		Body:     "body",
	}); err != nil {
		t.Fatalf("Send() error = %v", err)
	}
}

func TestClientSendRejectsNonSuccessStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	client := New(server.Client())
	err := client.Send(t.Context(), contract.WXPushSendInput{
		BaseURL:  server.URL,
		APIToken: "token",
		Title:    "title",
		Body:     "body",
	})
	if err == nil {
		t.Fatal("Send() expected error")
	}
}
