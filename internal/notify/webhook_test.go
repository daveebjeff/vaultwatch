package notify

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewWebhookNotifier_EmptyURL(t *testing.T) {
	_, err := NewWebhookNotifier("")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestWebhookNotifier_Send_Success(t *testing.T) {
	var got webhookPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &got)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, err := NewWebhookNotifier(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msg := Message{
		Path:      "secret/api-key",
		Status:    StatusExpiringSoon,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Path != msg.Path {
		t.Errorf("expected path %q, got %q", msg.Path, got.Path)
	}
	if got.Status != string(StatusExpiringSoon) {
		t.Errorf("expected status %q, got %q", StatusExpiringSoon, got.Status)
	}
}

func TestWebhookNotifier_Send_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n, _ := NewWebhookNotifier(ts.URL)
	msg := Message{
		Path:      "secret/api-key",
		Status:    StatusExpired,
		ExpiresAt: time.Now(),
	}
	if err := n.Send(msg); err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestWebhookNotifier_Send_ContentType(t *testing.T) {
	var contentType string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n, _ := NewWebhookNotifier(ts.URL)
	_ = n.Send(Message{Path: "secret/test", Status: StatusExpiringSoon, ExpiresAt: time.Now()})

	if contentType != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", contentType)
	}
}
