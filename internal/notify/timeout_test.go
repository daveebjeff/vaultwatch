package notify

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewTimeoutNotifier_NilInner(t *testing.T) {
	_, err := NewTimeoutNotifier(nil, time.Second)
	if err == nil {
		t.Fatal("expected error for nil inner notifier")
	}
}

func TestNewTimeoutNotifier_ZeroDuration(t *testing.T) {
	_, err := NewTimeoutNotifier(NewNoopNotifier(), 0)
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
}

func TestNewTimeoutNotifier_Valid(t *testing.T) {
	tn, err := NewTimeoutNotifier(NewNoopNotifier(), time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tn == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestTimeoutNotifier_Send_Success(t *testing.T) {
	n, _ := NewTimeoutNotifier(NewNoopNotifier(), time.Second)
	err := n.Send(context.Background(), exampleMsg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTimeoutNotifier_Send_ExceedsDeadline(t *testing.T) {
	slow := &mockSlowNotifier{delay: 200 * time.Millisecond}
	n, _ := NewTimeoutNotifier(slow, 30*time.Millisecond)
	err := n.Send(context.Background(), exampleMsg)
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got: %v", err)
	}
}

func TestTimeoutNotifier_Send_InnerError(t *testing.T) {
	sentinel := errors.New("inner failure")
	failing := &mockFailNotifier{err: sentinel}
	n, _ := NewTimeoutNotifier(failing, time.Second)
	err := n.Send(context.Background(), exampleMsg)
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got: %v", err)
	}
}

// mockSlowNotifier sleeps for delay before returning.
type mockSlowNotifier struct{ delay time.Duration }

func (m *mockSlowNotifier) Send(ctx context.Context, _ Message) error {
	select {
	case <-time.After(m.delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// mockFailNotifier always returns the configured error.
type mockFailNotifier struct{ err error }

func (m *mockFailNotifier) Send(_ context.Context, _ Message) error { return m.err }
