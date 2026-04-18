package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GoogleChatNotifier sends alerts to a Google Chat webhook.
type GoogleChatNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewGoogleChatNotifier creates a new GoogleChatNotifier.
func NewGoogleChatNotifier(webhookURL string) (*GoogleChatNotifier, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("googlechat: webhook URL must not be empty")
	}
	return &GoogleChatNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}, nil
}

// Send delivers a notification message to Google Chat.
func (g *GoogleChatNotifier) Send(msg Message) error {
	body := map[string]string{
		"text": fmt.Sprintf("*[%s] %s*\nSecret: %s\nExpires: %s",
			msg.Status, msg.Title, msg.SecretPath, msg.ExpiresAt.Format("2006-01-02 15:04:05 UTC")),
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("googlechat: failed to marshal payload: %w", err)
	}
	resp, err := g.client.Post(g.webhookURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("googlechat: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("googlechat: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
