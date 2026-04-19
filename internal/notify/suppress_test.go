package notify

import (
	"errors"
	"testing"
	"time"
)

func TestNewSuppressNotifier_NilInner(t *testing.T) {
	_, err := NewSuppressNotifier(nil, time.Minute)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewSuppressNotifier_ZeroTTL(t *testing.T) {
	_, err := NewSuppressNotifier(NewNoopNotifier(), 0)
	if err == nil {
		t.Fatal("expected error for zero ttl")
	}
}

func TestNewSuppressNotifier_Valid(t *testing.T) {
	s, err := NewSuppressNotifier(NewNoopNotifier(), time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestSuppressNotifier_ForwardsWhenNotSuppressed(t *testing.T) {
	called := 0
	mock := &mockNotifier{fn: func(msg Message) error { called++; return nil }}
	s, _ := NewSuppressNotifier(mock, time.Minute)

	_ = s.Send(Message{Path: "secret/foo", Status: StatusExpiringSoon})
	if called != 1 {
		t.Fatalf("expected 1 call, got %d", called)
	}
}

func TestSuppressNotifier_SuppressesPath(t *testing.T) {
	called := 0
	mock := &mockNotifier{fn: func(msg Message) error { called++; return nil }}
	s, _ := NewSuppressNotifier(mock, time.Minute)

	s.Suppress("secret/foo")
	_ = s.Send(Message{Path: "secret/foo", Status: StatusExpiringSoon})
	if called != 0 {
		t.Fatalf("expected 0 calls while suppressed, got %d", called)
	}
}

func TestSuppressNotifier_UnsuppressLifts(t *testing.T) {
	called := 0
	mock := &mockNotifier{fn: func(msg Message) error { called++; return nil }}
	s, _ := NewSuppressNotifier(mock, time.Minute)

	s.Suppress("secret/foo")
	s.Unsuppress("secret/foo")
	_ = s.Send(Message{Path: "secret/foo", Status: StatusExpiringSoon})
	if called != 1 {
		t.Fatalf("expected 1 call after unsuppress, got %d", called)
	}
}

func TestSuppressNotifier_ExpiredSuppressionForwards(t *testing.T) {
	called := 0
	mock := &mockNotifier{fn: func(msg Message) error { called++; return nil }}
	s, _ := NewSuppressNotifier(mock, time.Millisecond)

	s.Suppress("secret/bar")
	time.Sleep(5 * time.Millisecond)
	_ = s.Send(Message{Path: "secret/bar", Status: StatusExpired})
	if called != 1 {
		t.Fatalf("expected 1 call after ttl expiry, got %d", called)
	}
}

func TestSuppressNotifier_InnerErrorReturned(t *testing.T) {
	want := errors.New("send failed")
	mock := &mockNotifier{fn: func(msg Message) error { return want }}
	s, _ := NewSuppressNotifier(mock, time.Minute)

	err := s.Send(Message{Path: "secret/baz"})
	if !errors.Is(err, want) {
		t.Fatalf("expected wrapped error, got %v", err)
	}
}
