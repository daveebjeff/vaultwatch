package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const datadogEventsURL = "https://api.datadoghq.com/api/v1/events"

// DatadogNotifier sends alerts as Datadog events.
type DatadogNotifier struct {
	apiKey string
	url    string
	client *http.Client
}

// NewDatadogNotifier creates a new DatadogNotifier.
func NewDatadogNotifier(apiKey string) (*DatadogNotifier, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("datadog: API key must not be empty")
	}
	return &DatadogNotifier{
		apiKey: apiKey,
		url:    datadogEventsURL,
		client: &http.Client{},
	}, nil
}

type datadogEvent struct {
	Title     string   `json:"title"`
	Text      string   `json:"text"`
	AlertType string   `json:"alert_type"`
	Tags      []string `json:"tags"`
}

func datadogAlertType(s Status) string {
	switch s {
	case StatusExpired:
		return "error"
	case StatusExpiringSoon:
		return "warning"
	default:
		return "info"
	}
}

// Send delivers the message as a Datadog event.
func (d *DatadogNotifier) Send(msg Message) error {
	event := datadogEvent{
		Title:     fmt.Sprintf("VaultWatch: %s", msg.SecretPath),
		Text:      msg.Body,
		AlertType: datadogAlertType(msg.Status),
		Tags:      []string{"source:vaultwatch", fmt.Sprintf("secret:%s", msg.SecretPath)},
	}
	b, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("datadog: marshal event: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, d.url, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("datadog: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("DD-API-KEY", d.apiKey)
	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("datadog: http post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("datadog: unexpected status %d", resp.StatusCode)
	}
	return nil
}
