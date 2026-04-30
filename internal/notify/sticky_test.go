package notify

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNewStickyNotifier_NilInner(t *testing.T) {
	_, err := NewStickyNotifier(nil, time.Second)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewStickyNotifier_ZeroInterval(t *testing.T) {
	_, err := NewStickyNotifier(NewNoopNotifier(), 0)
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestNewStickyNotifier_Valid(t *testing.T) {
	sn, err := NewStickyNotifier(NewNoopNotifier(), 10*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sn == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestStickyNotifier_Send_ForwardsImmediately(t *testing.T) {
	var mu sync.Mutex
	var count int
	capture := &callCountNotifier{fn: func() { mu.Lock(); count++; mu.Unlock() }}

	sn, _ := NewStickyNotifier(capture, 20*time.Millisecond)
	msg := Message{Path: "secret/foo", Status: StatusExpiringSoon}

	if err := sn.Send(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	got := count
	mu.Unlock()
	if got != 1 {
		t.Fatalf("expected 1 immediate send, got %d", got)
	}

	sn.Clear(msg.Path)
}

func TestStickyNotifier_RepeatsUntilCleared(t *testing.T) {
	var mu sync.Mutex
	var count int
	capture := &callCountNotifier{fn: func() { mu.Lock(); count++; mu.Unlock() }}

	sn, _ := NewStickyNotifier(capture, 15*time.Millisecond)
	msg := Message{Path: "secret/bar", Status: StatusExpired}

	_ = sn.Send(context.Background(), msg)
	time.Sleep(55 * time.Millisecond)
	sn.Clear(msg.Path)

	mu.Lock()
	got := count
	mu.Unlock()
	// Should have fired at t=0, ~15ms, ~30ms, ~45ms → at least 3 total
	if got < 3 {
		t.Fatalf("expected at least 3 sends, got %d", got)
	}
}

func TestStickyNotifier_ClearStopsRepeat(t *testing.T) {
	var mu sync.Mutex
	var count int
	capture := &callCountNotifier{fn: func() { mu.Lock(); count++; mu.Unlock() }}

	sn, _ := NewStickyNotifier(capture, 10*time.Millisecond)
	msg := Message{Path: "secret/baz", Status: StatusExpiringSoon}

	_ = sn.Send(context.Background(), msg)
	sn.Clear(msg.Path)
	time.Sleep(40 * time.Millisecond)

	mu.Lock()
	got := count
	mu.Unlock()
	if got != 1 {
		t.Fatalf("expected exactly 1 send after immediate clear, got %d", got)
	}
}

func TestStickyNotifier_ActivePaths(t *testing.T) {
	sn, _ := NewStickyNotifier(NewNoopNotifier(), 50*time.Millisecond)
	_ = sn.Send(context.Background(), Message{Path: "a"})
	_ = sn.Send(context.Background(), Message{Path: "b"})

	paths := sn.ActivePaths()
	if len(paths) != 2 {
		t.Fatalf("expected 2 active paths, got %d", len(paths))
	}
	sn.Clear("a")
	sn.Clear("b")
}

// callCountNotifier is a test helper that increments a counter on each Send.
type callCountNotifier struct{ fn func() }

func (c *callCountNotifier) Send(_ context.Context, _ Message) error {
	c.fn()
	return nil
}
