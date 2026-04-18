package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// VictorOpsNotifier sends alerts to VictorOps (Splunk On-Call) via REST endpoint.
type VictorOpsNotifier struct {
	webhookURL string
	client     *http.Client
}

// NewVictorOpsNotifier creates a new VictorOpsNotifier.
func NewVictorOpsNotifier(webhookURL string) (*VictorOpsNotifier, error) {
	if webhookURL == "" {
		return nil, fmt.Errorf("victorops: webhook URL must not be empty")
	}
	return &VictorOpsNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}, nil
}

type victorOpsPayload struct {
	MessageType       string `json:"message_type"`
	EntityDisplayName string `json:"entity_display_name"`
	StateMessage      string `json:"state_message"`
}

func victorOpsMessageType(s Status) string {
	switch s {
	case StatusExpired:
		return "CRITICAL"
	case StatusExpiringSoon:
		return "WARNING"
	default:
		return "INFO"
	}
}

// Send delivers the message to VictorOps.
func (v *VictorOpsNotifier) Send(msg Message) error {
	payload := victorOpsPayload{
		MessageType:       victorOpsMessageType(msg.Status),
		EntityDisplayName: fmt.Sprintf("VaultWatch: %s", msg.SecretPath),
		StateMessage:      msg.Body,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("victorops: marshal payload: %w", err)
	}
	resp, err := v.client.Post(v.webhookURL, "application/json", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("victorops: http post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("victorops: unexpected status %d", resp.StatusCode)
	}
	return nil
}
