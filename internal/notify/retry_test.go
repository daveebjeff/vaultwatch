package notify

import (
	"errors"
	"testing"
	"time"
)

type countingNotifier struct {
	calls int
	failUntil int
	err error
}

func (c *countingNotifier) Send(msg Message) error {
	c.calls++
	if c.calls <= c.failUntil {
		return c.err
	}
	return nil
}

func TestRetryNotifier_SuccessFirstTry(t *testing.T) {
	n := &countingNotifier{}
	r := NewRetryNotifier(n, 3, 0)
	if err := r.Send(exampleMsg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.calls != 1 {
		t.Errorf("expected 1 call, got %d", n.calls)
	}
}

func TestRetryNotifier_RetriesOnFailure(t *testing.T) {
	n := &countingNotifier{failUntil: 2, err: errors.New("temp error")}
	r := NewRetryNotifier(n, 3, 0)
	if err := r.Send(exampleMsg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.calls != 3 {
		t.Errorf("expected 3 calls, got %d", n.calls)
	}
}

func TestRetryNotifier_ExhaustsRetries(t *testing.T) {
	n := &countingNotifier{failUntil: 10, err: errors.New("persistent error")}
	r := NewRetryNotifier(n, 3, 0)
	if err := r.Send(exampleMsg); err == nil {
		t.Fatal("expected error but got nil")
	}
	if n.calls != 3 {
		t.Errorf("expected 3 calls, got %d", n.calls)
	}
}

func TestRetryNotifier_ZeroAttempts(t *testing.T) {
	n := &countingNotifier{}
	r := NewRetryNotifier(n, 0, 0)
	if err := r.Send(exampleMsg); err == nil {
		t.Fatal("expected error with zero attempts")
	}
	if n.calls != 0 {
		t.Errorf("expected 0 calls, got %d", n.calls)
	}
}

func TestRetryNotifier_DelayBetweenRetries(t *testing.T) {
	n := &countingNotifier{failUntil: 1, err: errors.New("err")}
	r := NewRetryNotifier(n, 2, 10*time.Millisecond)
	start := time.Now()
	_ = r.Send(exampleMsg)
	if time.Since(start) < 10*time.Millisecond {
		t.Error("expected delay between retries")
	}
}
