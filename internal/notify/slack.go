package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/vaultwatch/internal/alert"
)

// SlackNotifier sends alert notifications to a Slack webhook.
type SlackNotifier struct {
	webhookURL string
	client     *http.Client
}

type slackPayload struct {
	Text string `json:"text"`
}

// NewSlackNotifier creates a SlackNotifier that posts to the given webhook URL.
func NewSlackNotifier(webhookURL string) (*SlackNotifier, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("slack webhook URL must not be empty")
	}
	return &SlackNotifier{
		webhookURL: webhookURL,
		client: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send posts the alert message to Slack.
func (s *SlackNotifier) Send(msg alert.Message) error {
	text := fmt.Sprintf("[%s] %s — path: %s (expires: %s)",
		msg.Level,
		msg.Summary,
		msg.Path,
		msg.ExpiresAt.Format(time.RFC3339),
	)

	payload, err := json.Marshal(slackPayload{Text: text})
	if err != nil {
		return fmt.Errorf("slack: marshal payload: %w", err)
	}

	resp, err := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("slack: http post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack: unexpected status %d", resp.StatusCode)
	}
	return nil
}
