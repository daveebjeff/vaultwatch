package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewPagerDutyNotifier_EmptyKey(t *testing.T) {
	_, err := NewPagerDutyNotifier("")
	if err == nil {
		t.Fatal("expected error for empty integration key")
	}
}

func TestNewPagerDutyNotifier_Valid(t *testing.T) {
	n, err := NewPagerDutyNotifier("test-key-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestPagerDutyNotifier_Send_Success(t *testing.T) {
	var received bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received = true
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected application/json content-type")
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer ts.Close()

	n, _ := NewPagerDutyNotifier("key")
	n.endpoint = ts.URL

	msg := Message{
		Path:      "secret/db",
		Status:    StatusExpiringSoon,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !received {
		t.Fatal("server never received request")
	}
}

func TestPagerDutyNotifier_Send_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	n, _ := NewPagerDutyNotifier("key")
	n.endpoint = ts.URL

	msg := Message{
		Path:      "secret/db",
		Status:    StatusExpired,
		ExpiresAt: time.Now(),
	}
	if err := n.Send(msg); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}
