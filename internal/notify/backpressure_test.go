package notify

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewBackpressureNotifier_NilInner(t *testing.T) {
	_, err := NewBackpressureNotifier(nil, 10)
	if !errors.Is(err, ErrBackpressureNilInner) {
		t.Fatalf("expected ErrBackpressureNilInner, got %v", err)
	}
}

func TestNewBackpressureNotifier_ZeroCapacity(t *testing.T) {
	_, err := NewBackpressureNotifier(NewNoopNotifier(), 0)
	if !errors.Is(err, ErrBackpressureZeroCapacity) {
		t.Fatalf("expected ErrBackpressureZeroCapacity, got %v", err)
	}
}

func TestNewBackpressureNotifier_Valid(t *testing.T) {
	n, err := NewBackpressureNotifier(NewNoopNotifier(), 8)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer n.Stop()
}

func TestBackpressureNotifier_MessagesDelivered(t *testing.T) {
	var count atomic.Int64
	inner := &mockNotifier{fn: func(_ context.Context, _ Message) error {
		count.Add(1)
		return nil
	}}
	n, err := NewBackpressureNotifier(inner, 16)
	if err != nil {
		t.Fatal(err)
	}

	const total = 5
	for i := 0; i < total; i++ {
		if err := n.Send(context.Background(), Message{Path: "secret/a", Status: StatusExpiringSoon}); err != nil {
			t.Fatalf("unexpected send error: %v", err)
		}
	}
	n.Stop()

	if got := count.Load(); got != total {
		t.Fatalf("expected %d deliveries, got %d", total, got)
	}
}

func TestBackpressureNotifier_QueueFull(t *testing.T) {
	// Block the inner notifier so the queue fills up.
	block := make(chan struct{})
	inner := &mockNotifier{fn: func(_ context.Context, _ Message) error {
		<-block
		return nil
	}}
	n, err := NewBackpressureNotifier(inner, 2)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		close(block)
		n.Stop()
	}()

	msg := Message{Path: "secret/b", Status: StatusExpired}
	// Fill the queue; the drain goroutine will be blocked.
	time.Sleep(10 * time.Millisecond) // let drain goroutine start and block
	_ = n.Send(context.Background(), msg)
	_ = n.Send(context.Background(), msg)
	_ = n.Send(context.Background(), msg)

	err = n.Send(context.Background(), msg)
	if !errors.Is(err, ErrBackpressureQueueFull) {
		t.Fatalf("expected ErrBackpressureQueueFull, got %v", err)
	}
}

func TestBackpressureNotifier_StopDrainsQueue(t *testing.T) {
	var count atomic.Int64
	inner := &mockNotifier{fn: func(_ context.Context, _ Message) error {
		count.Add(1)
		return nil
	}}
	n, err := NewBackpressureNotifier(inner, 64)
	if err != nil {
		t.Fatal(err)
	}
	const total = 20
	for i := 0; i < total; i++ {
		_ = n.Send(context.Background(), Message{Path: "secret/c"})
	}
	n.Stop()
	if got := count.Load(); got != total {
		t.Fatalf("expected all %d messages drained, got %d", total, got)
	}
}
