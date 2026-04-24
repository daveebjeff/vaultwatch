package notify

import (
	"testing"
	"time"
)

func TestNewWindowNotifier_NilInner(t *testing.T) {
	_, err := NewWindowNotifier(nil, 3, time.Second)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewWindowNotifier_ZeroMax(t *testing.T) {
	_, err := NewWindowNotifier(NewNoopNotifier(), 0, time.Second)
	if err == nil {
		t.Fatal("expected error for zero max")
	}
}

func TestNewWindowNotifier_ZeroWindow(t *testing.T) {
	_, err := NewWindowNotifier(NewNoopNotifier(), 3, 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewWindowNotifier_Valid(t *testing.T) {
	wn, err := NewWindowNotifier(NewNoopNotifier(), 3, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if wn == nil {
		t.Fatal("expected non-nil WindowNotifier")
	}
}

func TestWindowNotifier_AllowsUpToMax(t *testing.T) {
	noop := NewNoopNotifier()
	wn, _ := NewWindowNotifier(noop, 3, time.Second)

	msg := Message{Path: "secret/foo", Status: StatusExpiringSoon}
	for i := 0; i < 3; i++ {
		if err := wn.Send(msg); err != nil {
			t.Fatalf("send %d: unexpected error: %v", i+1, err)
		}
	}
}

func TestWindowNotifier_SuppressesAboveMax(t *testing.T) {
	noop := NewNoopNotifier()
	wn, _ := NewWindowNotifier(noop, 2, time.Second)

	msg := Message{Path: "secret/bar", Status: StatusExpired}
	wn.Send(msg) //nolint:errcheck
	wn.Send(msg) //nolint:errcheck

	if err := wn.Send(msg); err != ErrSuppressed {
		t.Fatalf("expected ErrSuppressed, got %v", err)
	}
}

func TestWindowNotifier_AllowsAgainAfterWindowExpires(t *testing.T) {
	noop := NewNoopNotifier()
	wn, _ := NewWindowNotifier(noop, 1, 50*time.Millisecond)

	msg := Message{Path: "secret/baz", Status: StatusExpiringSoon}
	if err := wn.Send(msg); err != nil {
		t.Fatalf("first send: %v", err)
	}
	if err := wn.Send(msg); err != ErrSuppressed {
		t.Fatalf("expected suppression before window expires, got %v", err)
	}

	time.Sleep(60 * time.Millisecond)

	if err := wn.Send(msg); err != nil {
		t.Fatalf("expected send to succeed after window reset, got %v", err)
	}
}
