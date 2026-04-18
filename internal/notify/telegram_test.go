package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewTelegramNotifier_EmptyToken(t *testing.T) {
	_, err := NewTelegramNotifier("", "12345")
	if err == nil {
		t.Fatal("expected error for empty bot token")
	}
}

func TestNewTelegramNotifier_EmptyChatID(t *testing.T) {
	_, err := NewTelegramNotifier("token123", "")
	if err == nil {
		t.Fatal("expected error for empty chat ID")
	}
}

func TestNewTelegramNotifier_Valid(t *testing.T) {
	n, err := NewTelegramNotifier("token123", "-100123456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestTelegramNotifier_Send_Success(t *testing.T) {
	var called bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, _ := NewTelegramNotifier("token123", "-100123456")
	n.apiBase = ts.URL

	err := n.Send(Message{
		Title:      "Lease expiring",
		SecretPath: "secret/myapp",
		Status:     StatusExpiringSoon,
		ExpiresAt:  time.Now().Add(2 * time.Hour),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected server to be called")
	}
}

func TestTelegramNotifier_Send_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()

	n, _ := NewTelegramNotifier("badtoken", "-100123456")
	n.apiBase = ts.URL

	err := n.Send(Message{
		Title:      "Secret expired",
		SecretPath: "secret/old",
		Status:     StatusExpired,
		ExpiresAt:  time.Now().Add(-1 * time.Hour),
	})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}
