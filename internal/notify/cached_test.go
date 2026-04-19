package notify

import (
	"errors"
	"testing"
	"time"
)

func TestNewCachedNotifier_NilInner(t *testing.T) {
	_, err := NewCachedNotifier(nil, time.Minute)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewCachedNotifier_ZeroTTL(t *testing.T) {
	n := NewNoopNotifier()
	_, err := NewCachedNotifier(n, 0)
	if err == nil {
		t.Fatal("expected error for zero TTL")
	}
}

func TestNewCachedNotifier_Valid(t *testing.T) {
	n := NewNoopNotifier()
	c, err := NewCachedNotifier(n, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestCachedNotifier_FirstSendForwarded(t *testing.T) {
	var calls int
	mock := &mockNotifier{fn: func(Message) error { calls++; return nil }}
	c, _ := NewCachedNotifier(mock, time.Minute)

	msg := Message{Path: "secret/a", Status: StatusExpiringSoon}
	if err := c.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestCachedNotifier_DuplicateSuppressed(t *testing.T) {
	var calls int
	mock := &mockNotifier{fn: func(Message) error { calls++; return nil }}
	c, _ := NewCachedNotifier(mock, time.Minute)

	msg := Message{Path: "secret/a", Status: StatusExpiringSoon}
	_ = c.Send(msg)
	_ = c.Send(msg)
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestCachedNotifier_DifferentStatusForwarded(t *testing.T) {
	var calls int
	mock := &mockNotifier{fn: func(Message) error { calls++; return nil }}
	c, _ := NewCachedNotifier(mock, time.Minute)

	_ = c.Send(Message{Path: "secret/a", Status: StatusExpiringSoon})
	_ = c.Send(Message{Path: "secret/a", Status: StatusExpired})
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

func TestCachedNotifier_Invalidate(t *testing.T) {
	var calls int
	mock := &mockNotifier{fn: func(Message) error { calls++; return nil }}
	c, _ := NewCachedNotifier(mock, time.Minute)

	msg := Message{Path: "secret/b", Status: StatusExpired}
	_ = c.Send(msg)
	c.Invalidate(msg.Path, msg.Status)
	_ = c.Send(msg)
	if calls != 2 {
		t.Fatalf("expected 2 calls after invalidate, got %d", calls)
	}
}

func TestCachedNotifier_InnerErrorPropagated(t *testing.T) {
	sentinel := errors.New("send failed")
	mock := &mockNotifier{fn: func(Message) error { return sentinel }}
	c, _ := NewCachedNotifier(mock, time.Minute)

	err := c.Send(Message{Path: "secret/c", Status: StatusExpired})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}
