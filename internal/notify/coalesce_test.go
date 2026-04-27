package notify

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestNewCoalesceNotifier_NilInner(t *testing.T) {
	_, err := NewCoalesceNotifier(nil, 50*time.Millisecond)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewCoalesceNotifier_ZeroWindow(t *testing.T) {
	_, err := NewCoalesceNotifier(NewNoopNotifier(), 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewCoalesceNotifier_Valid(t *testing.T) {
	n, err := NewCoalesceNotifier(NewNoopNotifier(), 10*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestCoalesceNotifier_OnlyLastMessageForwarded(t *testing.T) {
	var received []Message
	inner := &mockNotifier{sendFn: func(m Message) error {
		received = append(received, m)
		return nil
	}}

	n, _ := NewCoalesceNotifier(inner, 60*time.Millisecond)

	msg1 := Message{Path: "secret/a", Body: "first"}
	msg2 := Message{Path: "secret/a", Body: "second"}
	msg3 := Message{Path: "secret/a", Body: "third"}

	_ = n.Send(msg1)
	_ = n.Send(msg2)
	_ = n.Send(msg3)

	time.Sleep(120 * time.Millisecond)

	if len(received) != 1 {
		t.Fatalf("expected 1 message, got %d", len(received))
	}
	if received[0].Body != "third" {
		t.Errorf("expected body 'third', got %q", received[0].Body)
	}
}

func TestCoalesceNotifier_IndependentPaths(t *testing.T) {
	var count int64
	inner := &mockNotifier{sendFn: func(m Message) error {
		atomic.AddInt64(&count, 1)
		return nil
	}}

	n, _ := NewCoalesceNotifier(inner, 50*time.Millisecond)

	_ = n.Send(Message{Path: "secret/a", Body: "a"})
	_ = n.Send(Message{Path: "secret/b", Body: "b"})

	time.Sleep(120 * time.Millisecond)

	if v := atomic.LoadInt64(&count); v != 2 {
		t.Errorf("expected 2 deliveries for distinct paths, got %d", v)
	}
}

func TestCoalesceNotifier_TimerResetOnNewMessage(t *testing.T) {
	var delivered int64
	inner := &mockNotifier{sendFn: func(m Message) error {
		atomic.AddInt64(&delivered, 1)
		return nil
	}}

	n, _ := NewCoalesceNotifier(inner, 80*time.Millisecond)

	// Send a message, then send another before window expires.
	_ = n.Send(Message{Path: "secret/x", Body: "v1"})
	time.Sleep(40 * time.Millisecond)
	_ = n.Send(Message{Path: "secret/x", Body: "v2"})

	// At 40ms the first timer should have been cancelled; nothing delivered yet.
	if v := atomic.LoadInt64(&delivered); v != 0 {
		t.Errorf("expected 0 deliveries at 40ms, got %d", v)
	}

	time.Sleep(120 * time.Millisecond)

	if v := atomic.LoadInt64(&delivered); v != 1 {
		t.Errorf("expected 1 delivery after window, got %d", v)
	}
}
