package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// MattermostNotifier sends alerts to a Mattermost channel via incoming webhook.
type MattermostNotifier struct {
	webhookURL string
	channel    string
	client     *http.Client
}

// NewMattermostNotifier creates a new MattermostNotifier.
func NewMattermostNotifier(webhookURL, channel string) (*MattermostNotifier, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("mattermost: webhook URL must not be empty")
	}
	return &MattermostNotifier{
		webhookURL: webhookURL,
		channel:    channel,
		client:     &http.Client{},
	}, nil
}

// Send delivers a notification message to Mattermost.
func (m *MattermostNotifier) Send(msg Message) error {
	text := fmt.Sprintf("**%s**\n%s", msg.Subject, msg.Body)

	payload := map[string]string{
		"text": text,
	}
	if m.channel != "" {
		payload["channel"] = m.channel
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("mattermost: failed to marshal payload: %w", err)
	}

	resp, err := m.client.Post(m.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("mattermost: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("mattermost: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
