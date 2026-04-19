package notify

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewPrometheusNotifier_NilInner(t *testing.T) {
	_, err := NewPrometheusNotifier(nil)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewPrometheusNotifier_Valid(t *testing.T) {
	p, err := NewPrometheusNotifier(NewNoopNotifier())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestPrometheusNotifier_Send_CountsStatus(t *testing.T) {
	p, _ := NewPrometheusNotifier(NewNoopNotifier())
	msg := Message{Path: "secret/foo", Status: StatusExpired, ExpiresAt: time.Now()}
	if err := p.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.counts["expired"] != 1 {
		t.Errorf("expected expired count 1, got %d", p.counts["expired"])
	}
}

func TestPrometheusNotifier_Send_ErrorCounted(t *testing.T) {
	failing := &mockFailNotifier{err: errors.New("boom")}
	p, _ := NewPrometheusNotifier(failing)
	msg := Message{Path: "secret/bar", Status: StatusExpiringSoon, ExpiresAt: time.Now()}
	_ = p.Send(msg)
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.counts["send_error"] != 1 {
		t.Errorf("expected send_error count 1, got %d", p.counts["send_error"])
	}
}

func TestPrometheusNotifier_Handler_Output(t *testing.T) {
	p, _ := NewPrometheusNotifier(NewNoopNotifier())
	msg := Message{Path: "secret/x", Status: StatusExpiringSoon, ExpiresAt: time.Now()}
	_ = p.Send(msg)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/metrics", nil)
	p.Handler()(rr, req)

	body := rr.Body.String()
	if !strings.Contains(body, "vaultwatch_notify_total") {
		t.Errorf("expected metric name in output, got: %s", body)
	}
	if !strings.Contains(body, "expiring_soon") {
		t.Errorf("expected expiring_soon label in output, got: %s", body)
	}
}

type mockFailNotifier struct{ err error }

func (m *mockFailNotifier) Send(_ Message) error { return m.err }
