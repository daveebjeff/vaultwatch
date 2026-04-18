package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// DiscordNotifier sends alerts to a Discord channel via webhook.
type DiscordNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewDiscordNotifier creates a new DiscordNotifier.
func NewDiscordNotifier(webhookURL string) (*DiscordNotifier, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("discord: webhook URL must not be empty")
	}
	return &DiscordNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}, nil
}

// Send delivers a notification message to Discord.
func (d *DiscordNotifier) Send(msg Message) error {
	color := 3066993 // green
	if msg.Status == StatusExpired {
		color = 15158332 // red
	} else if msg.Status == StatusExpiringSoon {
		color = 16776960 // yellow
	}

	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title":       msg.Subject,
				"description": msg.Body,
				"color":       color,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("discord: failed to marshal payload: %w", err)
	}

	resp, err := d.client.Post(d.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("discord: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
