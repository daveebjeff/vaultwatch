package notify

import (
	"fmt"
	"net/http"
	"net/url"
)

// HTTPGetNotifier sends an alert by making an HTTP GET request to a URL,
// appending status and secret path as query parameters.
type HTTPGetNotifier struct {
	baseURL string
	client  *http.Client
}

// NewHTTPGetNotifier creates a new HTTPGetNotifier.
func NewHTTPGetNotifier(baseURL string) (*HTTPGetNotifier, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("httpget: base URL must not be empty")
	}
	return &HTTPGetNotifier{
		baseURL: baseURL,
		client:  &http.Client{},
	}, nil
}

// Send performs a GET request with alert details as query parameters.
func (h *HTTPGetNotifier) Send(msg Message) error {
	u, err := url.Parse(h.baseURL)
	if err != nil {
		return fmt.Errorf("httpget: invalid base URL: %w", err)
	}
	q := u.Query()
	q.Set("status", string(msg.Status))
	q.Set("secret", msg.SecretPath)
	q.Set("expires_at", msg.ExpiresAt.UTC().Format("2006-01-02T15:04:05Z"))
	u.RawQuery = q.Encode()

	resp, err := h.client.Get(u.String())
	if err != nil {
		return fmt.Errorf("httpget: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("httpget: unexpected status code %d", resp.StatusCode)
	}
	return nil
}
