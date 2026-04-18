package notify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/alert"
)

func testMsg() alert.Message {
	return alert.Message{
		Level:     alert.LevelWarn,
		Summary:   "secret expiring soon",
		Path:      "secret/my-app/db",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
}

func TestNewSlackNotifier_EmptyURL(t *testing.T) {
	_, err := NewSlackNotifier("")
	if err == nil {
		t.Fatal("expected error for empty webhook URL")
	}
}

func TestSlackNotifier_Send_Success(t *testing.T) {
	var received slackPayload
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n, err := NewSlackNotifier(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msg := testMsg()
	if err := n.Send(msg); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	if !strings.Contains(received.Text, msg.Path) {
		t.Errorf("expected payload to contain path %q, got: %s", msg.Path, received.Text)
	}
	if !strings.Contains(received.Text, string(msg.Level)) {
		t.Errorf("expected payload to contain level %q, got: %s", msg.Level, received.Text)
	}
}

func TestSlackNotifier_Send_Non2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n, _ := NewSlackNotifier(server.URL)
	if err := n.Send(testMsg()); err == nil {
		t.Fatal("expected error on non-2xx response")
	}
}
