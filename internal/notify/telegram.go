package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const telegramAPIBase = "https://api.telegram.org"

// TelegramNotifier sends alerts via the Telegram Bot API.
type TelegramNotifier struct {
	botToken string
	chatID   string
	apiBase  string
	client   *http.Client
}

// NewTelegramNotifier creates a new TelegramNotifier.
func NewTelegramNotifier(botToken, chatID string) (*TelegramNotifier, error) {
	if botToken == "" {
		return nil, fmt.Errorf("telegram: bot token must not be empty")
	}
	if chatID == "" {
		return nil, fmt.Errorf("telegram: chat ID must not be empty")
	}
	return &TelegramNotifier{
		botToken: botToken,
		chatID:   chatID,
		apiBase:  telegramAPIBase,
		client:   &http.Client{},
	}, nil
}

// Send delivers a notification message via Telegram.
func (t *TelegramNotifier) Send(msg Message) error {
	text := fmt.Sprintf("<b>[%s] %s</b>\nSecret: <code>%s</code>\nExpires: %s",
		msg.Status, msg.Title, msg.SecretPath, msg.ExpiresAt.Format("2006-01-02 15:04:05 UTC"))
	payload, err := json.Marshal(map[string]string{
		"chat_id":    t.chatID,
		"text":       text,
		"parse_mode": "HTML",
	})
	if err != nil {
		return fmt.Errorf("telegram: failed to marshal payload: %w", err)
	}
	url := fmt.Sprintf("%s/bot%s/sendMessage", t.apiBase, t.botToken)
	resp, err := t.client.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("telegram: request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("telegram: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
