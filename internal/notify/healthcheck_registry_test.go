package notify

import (
	"testing"
	"time"
)

func TestHealthRegistry_RegisterAndStatuses(t *testing.T) {
	reg := NewHealthRegistry()
	h, _ := NewHealthCheckNotifier(NewNoopNotifier(), "noop", time.Hour)
	defer h.Stop()

	if err := reg.Register(h); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	statuses := reg.Statuses()
	if len(statuses) != 1 {
		t.Fatalf("expected 1 status, got %d", len(statuses))
	}
	if statuses[0].Name != "noop" {
		t.Errorf("unexpected name: %q", statuses[0].Name)
	}
}

func TestHealthRegistry_DuplicateRegister(t *testing.T) {
	reg := NewHealthRegistry()
	h, _ := NewHealthCheckNotifier(NewNoopNotifier(), "noop", time.Hour)
	defer h.Stop()

	_ = reg.Register(h)
	if err := reg.Register(h); err == nil {
		t.Fatal("expected error on duplicate register")
	}
}

func TestHealthRegistry_AllHealthy_True(t *testing.T) {
	reg := NewHealthRegistry()
	h, _ := NewHealthCheckNotifier(NewNoopNotifier(), "noop", time.Hour)
	defer h.Stop()
	_ = reg.Register(h)

	if !reg.AllHealthy() {
		t.Error("expected all healthy")
	}
}

func TestHealthRegistry_AllHealthy_False(t *testing.T) {
	reg := NewHealthRegistry()
	failing := &mockFailNotifier{err: errTest}
	h, _ := NewHealthCheckNotifier(failing, "failing", time.Hour)
	defer h.Stop()
	_ = reg.Register(h)

	// trigger a failure so health is recorded
	_ = h.Send(testCtx, Message{Path: "x", Status: StatusExpired})

	if reg.AllHealthy() {
		t.Error("expected not all healthy")
	}
}

func TestHealthRegistry_StopAll(t *testing.T) {
	reg := NewHealthRegistry()
	h, _ := NewHealthCheckNotifier(NewNoopNotifier(), "noop", 10*time.Millisecond)
	_ = reg.Register(h)
	// Should not panic or block.
	reg.StopAll()
}
