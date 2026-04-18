package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewDiscordNotifier_EmptyURL(t *testing.T) {
	_, err := NewDiscordNotifier("")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestNewDiscordNotifier_Valid(t *testing.T) {
	n, err := NewDiscordNotifier("https://discord.com/api/webhooks/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestDiscordNotifier_Send_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json content type")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	n, _ := NewDiscordNotifier(ts.URL)
	msg := Message{
		Subject:   "Test Alert",
		Body:      "Secret is expiring soon",
		Status:    StatusExpiringSoon,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		Path:      "secret/myapp",
	}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDiscordNotifier_Send_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n, _ := NewDiscordNotifier(ts.URL)
	msg := Message{
		Subject: "Test",
		Body:    "body",
		Status:  StatusExpired,
	}
	if err := n.Send(msg); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}
