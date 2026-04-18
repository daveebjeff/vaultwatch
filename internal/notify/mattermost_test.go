package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewMattermostNotifier_EmptyURL(t *testing.T) {
	_, err := NewMattermostNotifier("", "")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestNewMattermostNotifier_Valid(t *testing.T) {
	n, err := NewMattermostNotifier("https://mattermost.example.com/hooks/abc", "alerts")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestMattermostNotifier_Send_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json content type")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, _ := NewMattermostNotifier(ts.URL, "vault-alerts")
	msg := Message{
		Subject:   "Vault Alert",
		Body:      "Secret expiring in 5 minutes",
		Status:    StatusExpiringSoon,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Path:      "secret/db",
	}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMattermostNotifier_Send_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	n, _ := NewMattermostNotifier(ts.URL, "")
	msg := Message{
		Subject: "Test",
		Body:    "body",
		Status:  StatusExpired,
	}
	if err := n.Send(msg); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}
