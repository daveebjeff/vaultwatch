package notify

import (
	"context"
	"testing"
	"time"
)

func TestNewMuteNotifier_NilInner(t *testing.T) {
	_, err := NewMuteNotifier(nil)
	if err == nil {
		t.Fatal("expected error for nil inner, got nil")
	}
}

func TestNewMuteNotifier_Valid(t *testing.T) {
	n := NewNoopNotifier()
	m, err := NewMuteNotifier(n)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil MuteNotifier")
	}
}

func TestMuteNotifier_UnmutedForwards(t *testing.T) {
	rec := &recordingNotifier{}
	m, _ := NewMuteNotifier(rec)

	msg := Message{Path: "secret/foo", Status: StatusExpiringSoon}
	if err := m.Send(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.count != 1 {
		t.Fatalf("expected 1 send, got %d", rec.count)
	}
}

func TestMuteNotifier_MutedSuppresses(t *testing.T) {
	rec := &recordingNotifier{}
	m, _ := NewMuteNotifier(rec)

	m.Mute(5 * time.Second)

	msg := Message{Path: "secret/foo", Status: StatusExpiringSoon}
	if err := m.Send(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error during mute: %v", err)
	}
	if rec.count != 0 {
		t.Fatalf("expected 0 sends while muted, got %d", rec.count)
	}
}

func TestMuteNotifier_UnmuteRestoresSends(t *testing.T) {
	rec := &recordingNotifier{}
	m, _ := NewMuteNotifier(rec)

	m.Mute(10 * time.Second)
	m.Unmute()

	msg := Message{Path: "secret/bar", Status: StatusExpired}
	if err := m.Send(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.count != 1 {
		t.Fatalf("expected 1 send after unmute, got %d", rec.count)
	}
}

func TestMuteNotifier_IsMuted(t *testing.T) {
	n := NewNoopNotifier()
	m, _ := NewMuteNotifier(n)

	if m.IsMuted() {
		t.Fatal("expected not muted initially")
	}
	m.Mute(5 * time.Second)
	if !m.IsMuted() {
		t.Fatal("expected muted after Mute()")
	}
	m.Unmute()
	if m.IsMuted() {
		t.Fatal("expected not muted after Unmute()")
	}
}

func TestMuteNotifier_MuteExtends(t *testing.T) {
	n := NewNoopNotifier()
	m, _ := NewMuteNotifier(n)

	m.Mute(1 * time.Second)
	m.Mute(10 * time.Second)

	// After two Mute calls the deadline should be ~10s from now
	m.mu.Lock()
	remaining := time.Until(m.muteUntil)
	m.mu.Unlock()

	if remaining < 8*time.Second {
		t.Fatalf("expected mute window extended to ~10s, got %v", remaining)
	}
}
