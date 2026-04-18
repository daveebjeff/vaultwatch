package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// TeamsNotifier sends alerts to a Microsoft Teams channel via incoming webhook.
type TeamsNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewTeamsNotifier creates a TeamsNotifier. webhookURL must be non-empty.
func NewTeamsNotifier(webhookURL string) (*TeamsNotifier, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("teams: webhook URL must not be empty")
	}
	return &TeamsNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}, nil
}

// Send posts a message card to the Teams channel.
func (t *TeamsNotifier) Send(msg Message) error {
	payload := map[string]interface{}{
		"@type":      "MessageCard",
		"@context":   "http://schema.org/extensions",
		"summary":    msg.Summary,
		"themeColor": themeColor(msg.Status),
		"sections": []map[string]interface{}{
			{
				"activityTitle":    msg.Summary,
				"activitySubtitle": fmt.Sprintf("Path: %s", msg.SecretPath),
				"facts": []map[string]string{
					{"name": "Status", "value": string(msg.Status)},
					{"name": "Expires At", "value": msg.ExpiresAt.String()},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("teams: marshal payload: %w", err)
	}

	resp, err := t.client.Post(t.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("teams: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("teams: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func themeColor(s Status) string {
	switch s {
	case StatusExpired:
		return "FF0000"
	case StatusExpiringSoon:
		return "FFA500"
	default:
		return "00FF00"
	}
}
