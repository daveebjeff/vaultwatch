package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewDatadogNotifier_EmptyKey(t *testing.T) {
	_, err := NewDatadogNotifier("")
	if err == nil {
		t.Fatal("expected error for empty API key")
	}
}

func TestNewDatadogNotifier_Valid(t *testing.T) {
	n, err := NewDatadogNotifier("abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestDatadogNotifier_Send_Success(t *testing.T) {
	var gotAPIKey string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAPIKey = r.Header.Get("DD-API-KEY")
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n, _ := NewDatadogNotifier("test-key")
	n.url = ts.URL

	msg := Message{
		SecretPath: "secret/api",
		Status:     StatusExpiringSoon,
		ExpireAt:   time.Now().Add(24 * 60 * 60 * 1000000000),
		Body:       "Expiring soon",
	}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotAPIKey != "test-key" {
		t.Errorf("expected DD-API-KEY test-key, got %s", gotAPIKey)
	}
}

func TestDatadogNotifier_Send_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n, _ := NewDatadogNotifier("test-key")
	n.url = ts.URL

	msg := Message{
		SecretPath: "secret/api",
		Status:     StatusExpired,
		ExpireAt:   time.Now(),
		Body:       "Expired",
	}
	if err := n.Send(msg); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestDatadogAlertType(t *testing.T) {
	cases := []struct {
		status   Status
		expected string
	}{
		{StatusExpired, "error"},
		{StatusExpiringSoon, "warning"},
		{StatusOK, "info"},
	}
	for _, c := range cases {
		got := datadogAlertType(c.status)
		if got != c.expected {
			t.Errorf("status %v: expected %s, got %s", c.status, c.expected, got)
		}
	}
}
