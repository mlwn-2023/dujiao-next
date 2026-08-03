package wxpush

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dujiao-next/internal/modules/notification/contract"
)

const defaultTimeout = 15 * time.Second

type Client struct {
	httpClient *http.Client
	now        func() time.Time
}

func New(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}
	return &Client{httpClient: httpClient, now: time.Now}
}

type sendRequest struct {
	Timestamp string   `json:"timestamp"`
	Sign      string   `json:"sign"`
	Groups    []string `json:"groups,omitempty"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
}

func (c *Client) Send(ctx context.Context, input contract.WXPushSendInput) error {
	if c == nil || c.httpClient == nil {
		return errors.New("wxpush client is not configured")
	}
	baseURL := strings.TrimRight(strings.TrimSpace(input.BaseURL), "/")
	token := strings.TrimSpace(input.APIToken)
	if baseURL == "" || token == "" {
		return errors.New("wxpush endpoint or api token is empty")
	}

	timestamp := fmt.Sprintf("%d", c.now().UnixMilli())
	mac := hmac.New(sha256.New, []byte(token))
	_, _ = mac.Write([]byte(timestamp + "\n" + token))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	payload, err := json.Marshal(sendRequest{
		Timestamp: timestamp,
		Sign:      signature,
		Groups:    append([]string(nil), input.Groups...),
		Title:     strings.TrimSpace(input.Title),
		Content:   strings.TrimSpace(input.Body),
	})
	if err != nil {
		return fmt.Errorf("encode wxpush request: %w", err)
	}

	endpoint := baseURL
	if !strings.HasSuffix(strings.ToLower(endpoint), "/wxsend") {
		endpoint += "/wxsend"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("create wxpush request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send wxpush request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("wxpush request failed with status %d", resp.StatusCode)
	}
	return nil
}
