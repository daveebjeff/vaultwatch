package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewSplunkNotifier_EmptyURL(t *testing.T) {
	_, err := NewSplunkNotifier("", "token")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestNewSplunkNotifier_EmptyToken(t *testing.T) {
	_, err := NewSplunkNotifier("http://splunk:8088", "")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestNewSplunkNotifier_Valid(t *testing.T) {
	n, err := NewSplunkNotifier("http://splunk:8088/services/collector/event", "mytoken")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestSplunkNotifier_Send_Success(t *testing.T) {
	var gotAuth string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, _ := NewSplunkNotifier(ts.URL, "tok123")
	msg := Message{
		Path:      "secret/db",
		Status:    StatusExpired,
		ExpiresAt: time.Now(),
		Body:      "expired",
	}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotAuth != "Splunk tok123" {
		t.Errorf("expected Authorization header 'Splunk tok123', got %q", gotAuth)
	}
}

func TestSplunkNotifier_Send_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n, _ := NewSplunkNotifier(ts.URL, "tok")
	err := n.Send(Message{ExpiresAt: time.Now()})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestSplunkNotifier_Send_BadURL(t *testing.T) {
	n, _ := NewSplunkNotifier("http://127.0.0.1:0/bad", "tok")
	err := n.Send(Message{ExpiresAt: time.Now()})
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
