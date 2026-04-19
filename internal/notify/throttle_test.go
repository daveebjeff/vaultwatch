package notify

import (
	"errors"
	"testing"
	"time"
)

func TestNewThrottleNotifier_NilInner(t *testing.T) {
	_, err := NewThrottleNotifier(nil, 5, time.Minute)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewThrottleNotifier_ZeroMax(t *testing.T) {
	_, err := NewThrottleNotifier(NewNoopNotifier(), 0, time.Minute)
	if err == nil {
		t.Fatal("expected error for zero maxCount")
	}
}

func TestNewThrottleNotifier_ZeroWindow(t *testing.T) {
	_, err := NewThrottleNotifier(NewNoopNotifier(), 5, 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestThrottleNotifier_AllowsUpToMax(t *testing.T) {
	var calls int
	inner := &callCountNotifier{fn: func() error { calls++; return nil }}
	th, _ := NewThrottleNotifier(inner, 3, time.Minute)

	msg := Message{Path: "secret/foo", Status: StatusExpiringSoon}
	for i := 0; i < 5; i++ {
		th.Send(msg)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestThrottleNotifier_ResetsAfterWindow(t *testing.T) {
	var calls int
	inner := &callCountNotifier{fn: func() error { calls++; return nil }}
	th, _ := NewThrottleNotifier(inner, 2, 50*time.Millisecond)

	msg := Message{Path: "secret/bar", Status: StatusExpired}
	th.Send(msg)
	th.Send(msg)
	th.Send(msg) // dropped

	time.Sleep(60 * time.Millisecond)
	th.Send(msg) // window reset

	if calls != 3 {
		t.Fatalf("expected 3 calls after reset, got %d", calls)
	}
}

func TestThrottleNotifier_PropagatesError(t *testing.T) {
	inner := &callCountNotifier{fn: func() error { return errors.New("send failed") }}
	th, _ := NewThrottleNotifier(inner, 5, time.Minute)

	err := th.Send(Message{Path: "secret/x", Status: StatusExpired})
	if err == nil {
		t.Fatal("expected error to propagate")
	}
}

type callCountNotifier struct {
	fn func() error
}

func (c *callCountNotifier) Send(_ Message) error {
	return c.fn()
}
