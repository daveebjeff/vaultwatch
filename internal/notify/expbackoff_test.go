package notify

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

var errBackoff = errors.New("send failed")

func TestNewExpBackoffNotifier_NilInner(t *testing.T) {
	_, err := NewExpBackoffNotifier(nil, 10*time.Millisecond, 100*time.Millisecond, 3)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewExpBackoffNotifier_ZeroInitDelay(t *testing.T) {
	_, err := NewExpBackoffNotifier(NewNoopNotifier(), 0, 100*time.Millisecond, 3)
	if err == nil {
		t.Fatal("expected error for zero initDelay")
	}
}

func TestNewExpBackoffNotifier_Valid(t *testing.T) {
	n, err := NewExpBackoffNotifier(NewNoopNotifier(), 10*time.Millisecond, 100*time.Millisecond, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestExpBackoffNotifier_SuccessFirstAttempt(t *testing.T) {
	var calls int32
	mock := &mockNotifier{fn: func(_ context.Context, _ Message) error {
		atomic.AddInt32(&calls, 1)
		return nil
	}}
	n, _ := NewExpBackoffNotifier(mock, 10*time.Millisecond, 100*time.Millisecond, 3)
	if err := n.Send(context.Background(), Message{Path: "secret/a"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if atomic.LoadInt32(&calls) != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestExpBackoffNotifier_RetriesOnFailure(t *testing.T) {
	var calls int32
	mock := &mockNotifier{fn: func(_ context.Context, _ Message) error {
		n := atomic.AddInt32(&calls, 1)
		if n < 3 {
			return errBackoff
		}
		return nil
	}}
	n, _ := NewExpBackoffNotifier(mock, 5*time.Millisecond, 20*time.Millisecond, 5)
	if err := n.Send(context.Background(), Message{Path: "secret/b"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if atomic.LoadInt32(&calls) != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestExpBackoffNotifier_ExhaustsAttempts(t *testing.T) {
	mock := &mockNotifier{fn: func(_ context.Context, _ Message) error {
		return errBackoff
	}}
	n, _ := NewExpBackoffNotifier(mock, 5*time.Millisecond, 10*time.Millisecond, 3)
	if err := n.Send(context.Background(), Message{Path: "secret/c"}); !errors.Is(err, errBackoff) {
		t.Fatalf("expected errBackoff, got %v", err)
	}
}

func TestExpBackoffNotifier_ContextCancelled(t *testing.T) {
	mock := &mockNotifier{fn: func(_ context.Context, _ Message) error {
		return errBackoff
	}}
	n, _ := NewExpBackoffNotifier(mock, 50*time.Millisecond, 200*time.Millisecond, 5)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := n.Send(ctx, Message{Path: "secret/d"})
	if err == nil {
		t.Fatal("expected context error")
	}
}

func TestBuildExpBackoffNotifier_Defaults(t *testing.T) {
	cfg := ExpBackoffConfig{} // zero values → defaults applied
	n, err := BuildExpBackoffNotifier(NewNoopNotifier(), cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.attempts != DefaultExpBackoffConfig().Attempts {
		t.Fatalf("expected default attempts %d, got %d", DefaultExpBackoffConfig().Attempts, n.attempts)
	}
}
