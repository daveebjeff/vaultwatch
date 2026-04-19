package notify

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewConditionalNotifier_NilInner(t *testing.T) {
	_, err := NewConditionalNotifier(nil, func(Message) bool { return true })
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewConditionalNotifier_NilPredicate(t *testing.T) {
	noop := NewNoopNotifier()
	_, err := NewConditionalNotifier(noop, nil)
	if err == nil {
		t.Fatal("expected error for nil predicate")
	}
}

func TestNewConditionalNotifier_Valid(t *testing.T) {
	noop := NewNoopNotifier()
	c, err := NewConditionalNotifier(noop, func(Message) bool { return true })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestConditionalNotifier_PredicateTrue_Forwards(t *testing.T) {
	var received *Message
	inner := &mockNotifier{sendFn: func(_ context.Context, m Message) error {
		received = &m
		return nil
	}}
	c, _ := NewConditionalNotifier(inner, func(Message) bool { return true })
	msg := Message{Path: "secret/foo", Status: StatusExpired, ExpiresAt: time.Now()}
	if err := c.Send(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received == nil {
		t.Fatal("expected message to be forwarded")
	}
}

func TestConditionalNotifier_PredicateFalse_Suppresses(t *testing.T) {
	called := false
	inner := &mockNotifier{sendFn: func(_ context.Context, _ Message) error {
		called = true
		return nil
	}}
	c, _ := NewConditionalNotifier(inner, func(Message) bool { return false })
	msg := Message{Path: "secret/bar", Status: StatusExpiringSoon, ExpiresAt: time.Now()}
	if err := c.Send(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("expected message to be suppressed")
	}
}

func TestConditionalNotifier_InnerError_Propagated(t *testing.T) {
	sentinel := errors.New("inner failure")
	inner := &mockNotifier{sendFn: func(_ context.Context, _ Message) error {
		return sentinel
	}}
	c, _ := NewConditionalNotifier(inner, func(Message) bool { return true })
	msg := Message{Path: "secret/baz", Status: StatusExpired, ExpiresAt: time.Now()}
	if err := c.Send(context.Background(), msg); !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}
