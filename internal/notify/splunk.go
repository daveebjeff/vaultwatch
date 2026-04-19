package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SplunkNotifier sends alerts to a Splunk HTTP Event Collector (HEC) endpoint.
type SplunkNotifier struct {
	url   string
	token string
	client *http.Client
}

// NewSplunkNotifier creates a new SplunkNotifier.
// url is the HEC endpoint (e.g. https://splunk:8088/services/collector/event).
// token is the HEC token.
func NewSplunkNotifier(url, token string) (*SplunkNotifier, error) {
	if url == "" {
		return nil, fmt.Errorf("splunk: url must not be empty")
	}
	if token == "" {
		return nil, fmt.Errorf("splunk: token must not be empty")
	}
	return &SplunkNotifier{
		url:    url,
		token:  token,
		client: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// Send delivers a notification message to Splunk HEC.
func (s *SplunkNotifier) Send(msg Message) error {
	payload := map[string]interface{}{
		"time": msg.ExpiresAt.Unix(),
		"event": map[string]interface{}{
			"path":       msg.Path,
			"status":     msg.Status,
			"expires_at": msg.ExpiresAt.Format(time.RFC3339),
			"message":    msg.Body,
		},
		"sourcetype": "vaultwatch",
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("splunk: marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, s.url, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("splunk: create request: %w", err)
	}
	req.Header.Set("Authorization", "Splunk "+s.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("splunk: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("splunk: unexpected status %d", resp.StatusCode)
	}
	return nil
}
