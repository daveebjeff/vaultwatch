package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewGoogleChatNotifier_EmptyURL(t *testing.T) {
	_, err := NewGoogleChatNotifier("")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestNewGoogleChatNotifier_Valid(t *testing.T) {
	n, err := NewGoogleChatNotifier("https://chat.googleapis.com/v1/spaces/xxx/messages?key=yyy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestGoogleChatNotifier_Send_Success(t *testing.T) {
	var called bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, _ := NewGoogleChatNotifier(ts.URL)
	err := n.Send(Message{
		Title:      "Secret expiring",
		SecretPath: "secret/db",
		Status:     StatusExpiringSoon,
		ExpiresAt:  time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected server to be called")
	}
}

func TestGoogleChatNotifier_Send_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n, _ := NewGoogleChatNotifier(ts.URL)
	err := n.Send(Message{
		Title:      "Secret expired",
		SecretPath: "secret/api",
		Status:     StatusExpired,
		ExpiresAt:  time.Now().Add(-1 * time.Hour),
	})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}
