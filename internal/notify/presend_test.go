package notify

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewPreSendNotifier_NilInner(t *testing.T) {
	_, err := NewPreSendNotifier(nil, func(_ context.Context, _ Message) error { return nil })
	if err == nil {
		t.Fatal("expected error for nil inner, got nil")
	}
}

func TestNewPreSendNotifier_NilHook(t *testing.T) {
	noop := NewNoopNotifier()
	_, err := NewPreSendNotifier(noop, nil)
	if err == nil {
		t.Fatal("expected error for nil hook, got nil")
	}
}

func TestNewPreSendNotifier_Valid(t *testing.T) {
	noop := NewNoopNotifier()
	n, err := NewPreSendNotifier(noop, func(_ context.Context, _ Message) error { return nil })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestPreSendNotifier_HookAllows(t *testing.T) {
	var received Message
	inner := &capturingNotifier{}
	hook := func(_ context.Context, _ Message) error { return nil }

	n, _ := NewPreSendNotifier(inner, hook)
	msg := Message{
		Path:      "secret/db",
		Status:    StatusExpiringSoon,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		Body:      "expiring soon",
	}
	if err := n.Send(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	received = inner.last
	if received.Path != msg.Path {
		t.Errorf("expected path %q, got %q", msg.Path, received.Path)
	}
}

func TestPreSendNotifier_HookRejects(t *testing.T) {
	inner := &capturingNotifier{}
	hookErr := errors.New("validation failed")
	hook := func(_ context.Context, _ Message) error { return hookErr }

	n, _ := NewPreSendNotifier(inner, hook)
	msg := Message{Path: "secret/token", Status: StatusExpired, Body: "expired"}

	err := n.Send(context.Background(), msg)
	if err == nil {
		t.Fatal("expected error from hook rejection, got nil")
	}
	if !errors.Is(err, hookErr) {
		t.Errorf("expected wrapped hookErr, got: %v", err)
	}
	if inner.calls != 0 {
		t.Errorf("inner should not have been called, got %d calls", inner.calls)
	}
}

func TestPreSendNotifier_InnerErrorPropagated(t *testing.T) {
	innerErr := errors.New("delivery failed")
	inner := &failingNotifier{err: innerErr}
	hook := func(_ context.Context, _ Message) error { return nil }

	n, _ := NewPreSendNotifier(inner, hook)
	msg := Message{Path: "secret/cert", Status: StatusExpiringSoon, Body: "soon"}

	err := n.Send(context.Background(), msg)
	if !errors.Is(err, innerErr) {
		t.Errorf("expected innerErr, got: %v", err)
	}
}

// capturingNotifier records the last message and call count.
type capturingNotifier struct {
	last  Message
	calls int
}

func (c *capturingNotifier) Send(_ context.Context, msg Message) error {
	c.last = msg
	c.calls++
	return nil
}

// failingNotifier always returns the configured error.
type failingNotifier struct{ err error }

func (f *failingNotifier) Send(_ context.Context, _ Message) error { return f.err }
