package notify

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewHealthCheckNotifier_NilInner(t *testing.T) {
	_, err := NewHealthCheckNotifier(nil, "test", time.Second)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewHealthCheckNotifier_ZeroInterval(t *testing.T) {
	_, err := NewHealthCheckNotifier(NewNoopNotifier(), "test", 0)
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestNewHealthCheckNotifier_Valid(t *testing.T) {
	h, err := NewHealthCheckNotifier(NewNoopNotifier(), "noop", 10*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer h.Stop()
	if h.Status().Name != "noop" {
		t.Errorf("expected name 'noop', got %q", h.Status().Name)
	}
}

func TestHealthCheckNotifier_Send_RecordsHealthy(t *testing.T) {
	h, _ := NewHealthCheckNotifier(NewNoopNotifier(), "noop", time.Hour)
	defer h.Stop()

	err := h.Send(context.Background(), Message{Path: "secret/foo", Status: StatusExpired})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !h.Status().Healthy {
		t.Error("expected healthy after successful send")
	}
}

func TestHealthCheckNotifier_Send_RecordsUnhealthy(t *testing.T) {
	failing := &mockFailNotifier{err: errors.New("boom")}
	h, _ := NewHealthCheckNotifier(failing, "failing", time.Hour)
	defer h.Stop()

	_ = h.Send(context.Background(), Message{Path: "secret/foo", Status: StatusExpired})
	if h.Status().Healthy {
		t.Error("expected unhealthy after failed send")
	}
	if h.Status().LastError == nil {
		t.Error("expected LastError to be set")
	}
}

func TestHealthCheckNotifier_BackgroundProbe(t *testing.T) {
	n := NewNoopNotifier()
	h, _ := NewHealthCheckNotifier(n, "probe", 20*time.Millisecond)
	defer h.Stop()

	time.Sleep(60 * time.Millisecond)
	s := h.Status()
	if !s.Healthy {
		t.Error("expected healthy after background probe")
	}
}

// mockFailNotifier always returns an error.
type mockFailNotifier struct{ err error }

func (m *mockFailNotifier) Send(_ context.Context, _ Message) error { return m.err }
