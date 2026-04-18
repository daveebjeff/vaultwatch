package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewOpsGenieNotifier_EmptyKey(t *testing.T) {
	_, err := NewOpsGenieNotifier("")
	if err == nil {
		t.Fatal("expected error for empty api key")
	}
}

func TestNewOpsGenieNotifier_Valid(t *testing.T) {
	n, err := NewOpsGenieNotifier("test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestOpsGenieNotifier_Send_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			t.Error("expected Authorization header")
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	n, _ := NewOpsGenieNotifier("test-key")
	n.apiURL = server.URL

	msg := Message{
		SecretPath: "secret/db",
		Status:     StatusExpired,
		Expiry:     time.Now(),
		Detail:     "secret has expired",
	}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOpsGenieNotifier_Send_Non2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	n, _ := NewOpsGenieNotifier("bad-key")
	n.apiURL = server.URL

	msg := Message{
		SecretPath: "secret/db",
		Status:     StatusExpiringSoon,
		Expiry:     time.Now().Add(1 * time.Hour),
		Detail:     "expiring soon",
	}
	if err := n.Send(msg); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}
