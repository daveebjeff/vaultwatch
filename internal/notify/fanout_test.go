package notify

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

type countingNotifier struct {
	calls atomic.Int32
	err   error
}

func (c *countingNotifier) Send(_ context.Context, _ Message) error {
	c.calls.Add(1)
	return c.err
}

func TestNewFanoutNotifier_NoNotifiers(t *testing.T) {
	_, err := NewFanoutNotifier()
	if err == nil {
		t.Fatal("expected error for empty notifier list")
	}
}

func TestNewFanoutNotifier_NilNotifier(t *testing.T) {
	_, err := NewFanoutNotifier(nil)
	if err == nil {
		t.Fatal("expected error for nil notifier")
	}
}

func TestNewFanoutNotifier_Valid(t *testing.T) {
	_, err := NewFanoutNotifier(&countingNotifier{}, &countingNotifier{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFanoutNotifier_AllCalled(t *testing.T) {
	a := &countingNotifier{}
	b := &countingNotifier{}
	f, _ := NewFanoutNotifier(a, b)

	msg := Message{Path: "secret/foo", Status: StatusExpiringSoon, ExpiresAt: time.Now().Add(time.Hour)}
	if err := f.Send(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.calls.Load() != 1 || b.calls.Load() != 1 {
		t.Errorf("expected each notifier called once, got a=%d b=%d", a.calls.Load(), b.calls.Load())
	}
}

func TestFanoutNotifier_CollectsErrors(t *testing.T) {
	a := &countingNotifier{err: errors.New("boom")}
	b := &countingNotifier{err: errors.New("bang")}
	f, _ := NewFanoutNotifier(a, b)

	err := f.Send(context.Background(), Message{Path: "secret/bar"})
	if err == nil {
		t.Fatal("expected combined error")
	}
	if a.calls.Load() != 1 || b.calls.Load() != 1 {
		t.Error("all notifiers should be attempted even on error")
	}
}

func TestFanoutNotifier_PartialError(t *testing.T) {
	ok := &countingNotifier{}
	bad := &countingNotifier{err: errors.New("oops")}
	f, _ := NewFanoutNotifier(ok, bad)

	err := f.Send(context.Background(), Message{Path: "secret/baz"})
	if err == nil {
		t.Fatal("expected error from failing notifier")
	}
	if ok.calls.Load() != 1 {
		t.Error("successful notifier should still be called")
	}
}
