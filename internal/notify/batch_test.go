package notify

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewBatchNotifier_NilInner(t *testing.T) {
	_, err := NewBatchNotifier(nil, time.Second, 10)
	if !errors.Is(err, ErrNilInner) {
		t.Fatalf("expected ErrNilInner, got %v", err)
	}
}

func TestNewBatchNotifier_ZeroWindow(t *testing.T) {
	_, err := NewBatchNotifier(NewNoopNotifier(), 0, 10)
	if !errors.Is(err, ErrZeroWindow) {
		t.Fatalf("expected ErrZeroWindow, got %v", err)
	}
}

func TestNewBatchNotifier_Valid(t *testing.T) {
	b, err := NewBatchNotifier(NewNoopNotifier(), time.Second, 5)
	if err != nil {
		t.Fatal(err)
	}
	if b == nil {
		t.Fatal("expected non-nil BatchNotifier")
	}
}

func TestBatchNotifier_FlushOnMaxSize(t *testing.T) {
	var sent []Message
	inner := &recordNotifier{fn: func(m Message) { sent = append(sent, m) }}
	b, _ := NewBatchNotifier(inner, 10*time.Second, 3)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		_ = b.Send(ctx, Message{Path: "secret/a", Summary: "alert"})
	}
	if len(sent) != 1 {
		t.Fatalf("expected 1 flush, got %d", len(sent))
	}
}

func TestBatchNotifier_ManualFlush(t *testing.T) {
	var sent []Message
	inner := &recordNotifier{fn: func(m Message) { sent = append(sent, m) }}
	b, _ := NewBatchNotifier(inner, 10*time.Second, 10)
	ctx := context.Background()

	_ = b.Send(ctx, Message{Path: "secret/b", Summary: "one"})
	_ = b.Send(ctx, Message{Path: "secret/b", Summary: "two"})
	if err := b.Flush(ctx); err != nil {
		t.Fatal(err)
	}
	if len(sent) != 1 {
		t.Fatalf("expected 1 flushed message, got %d", len(sent))
	}
}

func TestBatchNotifier_EmptyFlush(t *testing.T) {
	noop := NewNoopNotifier()
	b, _ := NewBatchNotifier(noop, time.Second, 5)
	if err := b.Flush(context.Background()); err != nil {
		t.Fatalf("flush of empty batch should not error: %v", err)
	}
}

func TestBatchNotifier_SummaryMultiple(t *testing.T) {
	var got Message
	inner := &recordNotifier{fn: func(m Message) { got = m }}
	b, _ := NewBatchNotifier(inner, 10*time.Second, 2)
	ctx := context.Background()
	_ = b.Send(ctx, Message{Path: "p", Summary: "first"})
	_ = b.Send(ctx, Message{Path: "p", Summary: "second"})
	if got.Summary == "" {
		t.Fatal("expected non-empty summary")
	}
}

// recordNotifier is a test helper that records the last sent message.
type recordNotifier struct {
	fn func(Message)
}

func (r *recordNotifier) Send(_ context.Context, m Message) error {
	r.fn(m)
	return nil
}
