package notify

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewJitterNotifier_NilInner(t *testing.T) {
	_, err := NewJitterNotifier(nil, 10*time.Millisecond)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewJitterNotifier_ZeroMaxJitter(t *testing.T) {
	noop := NewNoopNotifier()
	_, err := NewJitterNotifier(noop, 0)
	if err == nil {
		t.Fatal("expected error for zero maxJitter")
	}
}

func TestNewJitterNotifier_NegativeMaxJitter(t *testing.T) {
	noop := NewNoopNotifier()
	_, err := NewJitterNotifier(noop, -1*time.Millisecond)
	if err == nil {
		t.Fatal("expected error for negative maxJitter")
	}
}

func TestNewJitterNotifier_Valid(t *testing.T) {
	noop := NewNoopNotifier()
	j, err := NewJitterNotifier(noop, 50*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if j == nil {
		t.Fatal("expected non-nil JitterNotifier")
	}
}

func TestJitterNotifier_Send_Success(t *testing.T) {
	noop := NewNoopNotifier()
	j, _ := NewJitterNotifier(noop, 5*time.Millisecond)
	msg := Message{Path: "secret/foo", Status: StatusExpiringSoon, ExpiresAt: time.Now().Add(time.Hour)}
	ctx := context.Background()
	if err := j.Send(ctx, msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestJitterNotifier_Send_ContextCancelled(t *testing.T) {
	// Use a large jitter so the context cancel fires first.
	noop := NewNoopNotifier()
	j, _ := NewJitterNotifier(noop, 10*time.Second)
	msg := Message{Path: "secret/bar", Status: StatusExpired, ExpiresAt: time.Now()}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	err := j.Send(ctx, msg)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
}

func TestJitterNotifier_Send_ForwardsInnerError(t *testing.T) {
	sentinel := errors.New("inner failure")
	failing := &mockNotifier{err: sentinel}
	j, _ := NewJitterNotifier(failing, 2*time.Millisecond)
	msg := Message{Path: "secret/baz", Status: StatusExpiringSoon}
	if err := j.Send(context.Background(), msg); !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}
