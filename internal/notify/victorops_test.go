package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewVictorOpsNotifier_EmptyURL(t *testing.T) {
	_, err := NewVictorOpsNotifier("")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestNewVictorOpsNotifier_Valid(t *testing.T) {
	n, err := NewVictorOpsNotifier("https://alert.victorops.com/integrations/generic/abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestVictorOpsNotifier_Send_Success(t *testing.T) {
	var gotContentType string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, _ := NewVictorOpsNotifier(ts.URL)
	msg := Message{
		SecretPath: "secret/db",
		Status:     StatusExpired,
		ExpireAt:   time.Now(),
		Body:       "Secret has expired",
	}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotContentType != "application/json" {
		t.Errorf("expected application/json, got %s", gotContentType)
	}
}

func TestVictorOpsNotifier_Send_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n, _ := NewVictorOpsNotifier(ts.URL)
	msg := Message{
		SecretPath: "secret/db",
		Status:     StatusExpiringSoon,
		ExpireAt:   time.Now(),
		Body:       "Expiring soon",
	}
	if err := n.Send(msg); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestVictorOpsMessageType(t *testing.T) {
	cases := []struct {
		status   Status
		expected string
	}{
		{StatusExpired, "CRITICAL"},
		{StatusExpiringSoon, "WARNING"},
		{StatusOK, "INFO"},
	}
	for _, c := range cases {
		got := victorOpsMessageType(c.status)
		if got != c.expected {
			t.Errorf("status %v: expected %s, got %s", c.status, c.expected, got)
		}
	}
}
