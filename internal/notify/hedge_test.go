package notify

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewHedgeNotifier_NilPrimary(t *testing.T) {
	_, err := NewHedgeNotifier(nil, NewNoopNotifier(), 10*time.Millisecond)
	if !errors.Is(err, ErrNilInner) {
		t.Fatalf("expected ErrNilInner, got %v", err)
	}
}

func TestNewHedgeNotifier_NilSecondary(t *testing.T) {
	_, err := NewHedgeNotifier(NewNoopNotifier(), nil, 10*time.Millisecond)
	if !errors.Is(err, errNilSecondary) {
		t.Fatalf("expected errNilSecondary, got %v", err)
	}
}

func TestNewHedgeNotifier_ZeroDelay(t *testing.T) {
	_, err := NewHedgeNotifier(NewNoopNotifier(), NewNoopNotifier(), 0)
	if !errors.Is(err, errZeroDuration) {
		t.Fatalf("expected errZeroDuration, got %v", err)
	}
}

func TestNewHedgeNotifier_Valid(t *testing.T) {
	h, err := NewHedgeNotifier(NewNoopNotifier(), NewNoopNotifier(), 10*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil HedgeNotifier")
	}
}

func TestHedgeNotifier_PrimarySucceedsFast(t *testing.T) {
	var secondaryCalled atomic.Bool
	secondary := &mockNotifier{fn: func(_ context.Context, _ Message) error {
		secondary Called.Store(true)
		return nil
	}}
	h, _ := NewHedgeNotifier(NewNoopNotifier(), secondary, 200*time.Millisecond)
	if err := h.Send(context.Background(), Message{Path: "secret/fast"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	time.Sleep(300 * time.Millisecond)
	if secondaryCalled.Load() {
		t.Error("secondary should not have been called when primary was fast")
	}
}

func TestHedgeNotifier_SecondaryCalledWhenPrimarySlow(t *testing.T) {
	var secondaryCalled atomic.Bool
	primary := &mockNotifier{fn: func(ctx context.Context, _ Message) error {
		select {
		case <-time.After(500 * time.Millisecond):
		case <-ctx.Done():
		}
		return errors.New("slow")
	}}
	secondary := &mockNotifier{fn: func(_ context.Context, _ Message) error {
		secondary Called.Store(true)
		return nil
	}}
	h, _ := NewHedgeNotifier(primary, secondary, 50*time.Millisecond)
	if err := h.Send(context.Background(), Message{Path: "secret/slow"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !secondaryCalled.Load() {
		t.Error("expected secondary to be called after hedge delay")
	}
}

func TestHedgeNotifier_BothFail_ReturnsPrimaryError(t *testing.T) {
	primaryErr := errors.New("primary failed")
	primary := &mockNotifier{fn: func(_ context.Context, _ Message) error { return primaryErr }}
	secondary := &mockNotifier{fn: func(_ context.Context, _ Message) error { return errors.New("secondary failed") }}
	h, _ := NewHedgeNotifier(primary, secondary, 1*time.Millisecond)
	err := h.Send(context.Background(), Message{Path: "secret/both-fail"})
	if !errors.Is(err, primaryErr) {
		t.Fatalf("expected primaryErr, got %v", err)
	}
}
