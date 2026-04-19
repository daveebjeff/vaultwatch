package notify

import (
	"errors"
	"testing"
	"time"
)

func TestNewAcknowledgeNotifier_NilInner(t *testing.T) {
	_, err := NewAcknowledgeNotifier(nil, time.Minute)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewAcknowledgeNotifier_ZeroTTL(t *testing.T) {
	_, err := NewAcknowledgeNotifier(NewNoopNotifier(), 0)
	if err == nil {
		t.Fatal("expected error for zero ttl")
	}
}

func TestNewAcknowledgeNotifier_Valid(t *testing.T) {
	n, err := NewAcknowledgeNotifier(NewNoopNotifier(), time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestAcknowledgeNotifier_ForwardsWhenNotAcked(t *testing.T) {
	var called int
	mock := &mockNotifier{fn: func(msg Message) error { called++; return nil }}
	n, _ := NewAcknowledgeNotifier(mock, time.Minute)

	msg := Message{Path: "secret/db"}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called != 1 {
		t.Fatalf("expected 1 call, got %d", called)
	}
}

func TestAcknowledgeNotifier_SuppressesWhenAcked(t *testing.T) {
	var called int
	mock := &mockNotifier{fn: func(msg Message) error { called++; return nil }}
	n, _ := NewAcknowledgeNotifier(mock, time.Minute)

	n.Acknowledge("secret/db")
	msg := Message{Path: "secret/db"}
	_ = n.Send(msg)
	if called != 0 {
		t.Fatalf("expected 0 calls while acked, got %d", called)
	}
}

func TestAcknowledgeNotifier_ResumeAfterExpiry(t *testing.T) {
	var called int
	mock := &mockNotifier{fn: func(msg Message) error { called++; return nil }}
	n, _ := NewAcknowledgeNotifier(mock, 10*time.Millisecond)

	n.Acknowledge("secret/db")
	time.Sleep(20 * time.Millisecond)

	msg := Message{Path: "secret/db"}
	_ = n.Send(msg)
	if called != 1 {
		t.Fatalf("expected 1 call after expiry, got %d", called)
	}
}

func TestAcknowledgeNotifier_DifferentPathsIndependent(t *testing.T) {
	var called int
	mock := &mockNotifier{fn: func(msg Message) error { called++; return nil }}
	n, _ := NewAcknowledgeNotifier(mock, time.Minute)

	n.Acknowledge("secret/db")
	_ = n.Send(Message{Path: "secret/db"})
	_ = n.Send(Message{Path: "secret/api"})
	if called != 1 {
		t.Fatalf("expected 1 call for unacked path, got %d", called)
	}
}

func TestAcknowledgeNotifier_PropagatesError(t *testing.T) {
	mock := &mockNotifier{fn: func(msg Message) error { return errors.New("send failed") }}
	n, _ := NewAcknowledgeNotifier(mock, time.Minute)
	err := n.Send(Message{Path: "secret/x"})
	if err == nil {
		t.Fatal("expected error to propagate")
	}
}
