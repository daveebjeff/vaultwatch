package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const opsgenieAPIURL = "https://api.opsgenie.com/v2/alerts"

// OpsGenieNotifier sends alerts to OpsGenie.
type OpsGenieNotifier struct {
	apiKey  string
	client  *http.Client
	apiURL  string
}

// NewOpsGenieNotifier creates a new OpsGenieNotifier.
func NewOpsGenieNotifier(apiKey string) (*OpsGenieNotifier, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("opsgenie: api key must not be empty")
	}
	return &OpsGenieNotifier{
		apiKey: apiKey,
		client: &http.Client{Timeout: 10 * time.Second},
		apiURL: opsgenieAPIURL,
	}, nil
}

type opsgeniePayload struct {
	Message     string `json:"message"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

// Send delivers a notification message to OpsGenie.
func (o *OpsGenieNotifier) Send(msg Message) error {
	priority := "P3"
	if msg.Status == StatusExpired {
		priority = "P1"
	}

	payload := opsgeniePayload{
		Message:     fmt.Sprintf("[VaultWatch] %s: %s", msg.Status, msg.SecretPath),
		Description: msg.Detail,
		Priority:    priority,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("opsgenie: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, o.apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("opsgenie: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "GenieKey "+o.apiKey)

	resp, err := o.client.Do(req)
	if err != nil {
		return fmt.Errorf("opsgenie: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("opsgenie: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
