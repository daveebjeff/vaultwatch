package notify

import (
	"context"
	"testing"
	"time"
)

func TestNewQuotaNotifier_NilInner(t *testing.T) {
	_, err := NewQuotaNotifier(nil, 10, time.Minute)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewQuotaNotifier_ZeroMax(t *testing.T) {
	_, err := NewQuotaNotifier(NewNoopNotifier(), 0, time.Minute)
	if err == nil {
		t.Fatal("expected error for zero max")
	}
}

func TestNewQuotaNotifier_ZeroWindow(t *testing.T) {
	_, err := NewQuotaNotifier(NewNoopNotifier(), 5, 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewQuotaNotifier_Valid(t *testing.T) {
	q, err := NewQuotaNotifier(NewNoopNotifier(), 5, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestQuotaNotifier_AllowsUpToMax(t *testing.T) {
	const max = 3
	q, _ := NewQuotaNotifier(NewNoopNotifier(), max, time.Hour)
	msg := Message{Path: "secret/test"}
	ctx := context.Background()

	for i := 0; i < max; i++ {
		if err := q.Send(ctx, msg); err != nil {
			t.Fatalf("send %d: unexpected error: %v", i+1, err)
		}
	}
}

func TestQuotaNotifier_ExceedsMax(t *testing.T) {
	const max = 2
	q, _ := NewQuotaNotifier(NewNoopNotifier(), max, time.Hour)
	msg := Message{Path: "secret/test"}
	ctx := context.Background()

	for i := 0; i < max; i++ {
		_ = q.Send(ctx, msg)
	}
	if err := q.Send(ctx, msg); err != ErrQuotaExceeded {
		t.Fatalf("expected ErrQuotaExceeded, got %v", err)
	}
}

func TestQuotaNotifier_ResetsAfterWindow(t *testing.T) {
	q, _ := NewQuotaNotifier(NewNoopNotifier(), 1, 20*time.Millisecond)
	msg := Message{Path: "secret/test"}
	ctx := context.Background()

	if err := q.Send(ctx, msg); err != nil {
		t.Fatalf("first send failed: %v", err)
	}
	if err := q.Send(ctx, msg); err != ErrQuotaExceeded {
		t.Fatalf("expected quota exceeded before reset, got %v", err)
	}

	time.Sleep(30 * time.Millisecond)

	if err := q.Send(ctx, msg); err != nil {
		t.Fatalf("send after window reset failed: %v", err)
	}
}

func TestQuotaNotifier_Remaining(t *testing.T) {
	const max = 5
	q, _ := NewQuotaNotifier(NewNoopNotifier(), max, time.Hour)
	msg := Message{Path: "secret/test"}
	ctx := context.Background()

	remaining, _ := q.Remaining()
	if remaining != max {
		t.Fatalf("expected %d remaining, got %d", max, remaining)
	}

	_ = q.Send(ctx, msg)
	_ = q.Send(ctx, msg)

	remaining, _ = q.Remaining()
	if remaining != max-2 {
		t.Fatalf("expected %d remaining, got %d", max-2, remaining)
	}
}
