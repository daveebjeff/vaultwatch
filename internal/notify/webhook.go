package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookNotifier posts a JSON payload to a configurable HTTP endpoint.
type WebhookNotifier struct {
	url    string
	client *http.Client
}

type webhookPayload struct {
	Path      string `json:"path"`
	Status    string `json:"status"`
	ExpiresAt string `json:"expires_at"`
	Message   string `json:"message"`
}

// NewWebhookNotifier creates a WebhookNotifier. url must not be empty.
func NewWebhookNotifier(url string) (*WebhookNotifier, error) {
	if url == "" {
		return nil, fmt.Errorf("webhook: url must not be empty")
	}
	return &WebhookNotifier{
		url:    url,
		client: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send posts the alert message as JSON to the configured webhook URL.
func (w *WebhookNotifier) Send(msg Message) error {
	payload := webhookPayload{
		Path:      msg.Path,
		Status:    string(msg.Status),
		ExpiresAt: msg.ExpiresAt.UTC().Format(time.RFC3339),
		Message:   fmt.Sprintf("[vaultwatch] secret %s is %s", msg.Path, msg.Status),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal: %w", err)
	}
	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("webhook: request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}
