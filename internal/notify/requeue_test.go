package notify

import (
	"errors"
	"testing"
	"time"
)

func TestNewRequeueNotifier_NilInner(t *testing.T) {
	_, err := NewRequeueNotifier(nil, 10, time.Second)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewRequeueNotifier_ZeroMax(t *testing.T) {
	_, err := NewRequeueNotifier(NewNoopNotifier(), 0, time.Second)
	if err == nil {
		t.Fatal("expected error for zero maxQueue")
	}
}

func TestNewRequeueNotifier_ZeroRetryAge(t *testing.T) {
	_, err := NewRequeueNotifier(NewNoopNotifier(), 10, 0)
	if err == nil {
		t.Fatal("expected error for zero retryAge")
	}
}

func TestNewRequeueNotifier_Valid(t *testing.T) {
	r, err := NewRequeueNotifier(NewNoopNotifier(), 10, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestRequeueNotifier_SuccessNotQueued(t *testing.T) {
	r, _ := NewRequeueNotifier(NewNoopNotifier(), 5, time.Minute)
	msg := Message{Path: "secret/a", Status: StatusExpiringSoon, At: time.Now()}
	if err := r.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.QueueLen() != 0 {
		t.Fatalf("expected empty queue, got %d", r.QueueLen())
	}
}

func TestRequeueNotifier_FailureQueued(t *testing.T) {
	fail := &mockNotifier{err: errors.New("send failed")}
	r, _ := NewRequeueNotifier(fail, 5, time.Minute)
	msg := Message{Path: "secret/b", Status: StatusExpired, At: time.Now()}
	_ = r.Send(msg)
	if r.QueueLen() != 1 {
		t.Fatalf("expected 1 queued message, got %d", r.QueueLen())
	}
}

func TestRequeueNotifier_MaxQueueDropsOldest(t *testing.T) {
	fail := &mockNotifier{err: errors.New("fail")}
	r, _ := NewRequeueNotifier(fail, 2, time.Minute)
	for i := 0; i < 4; i++ {
		_ = r.Send(Message{Path: "secret/x", At: time.Now()})
	}
	if r.QueueLen() > 2 {
		t.Fatalf("queue exceeded max: %d", r.QueueLen())
	}
}

func TestRequeueNotifier_FlushDelivers(t *testing.T) {
	fail := &mockNotifier{err: errors.New("fail")}
	r, _ := NewRequeueNotifier(fail, 5, time.Minute)
	_ = r.Send(Message{Path: "secret/c", At: time.Now()})
	if r.QueueLen() != 1 {
		t.Fatal("expected 1 queued message before flush")
	}
	// Allow delivery on flush.
	fail.err = nil
	r.Flush()
	if r.QueueLen() != 0 {
		t.Fatalf("expected empty queue after flush, got %d", r.QueueLen())
	}
}

func TestRequeueNotifier_AgedMessagesRetriedOnSend(t *testing.T) {
	fail := &mockNotifier{err: errors.New("fail")}
	// Use a very short retryAge so messages are eligible immediately.
	r, _ := NewRequeueNotifier(fail, 5, time.Nanosecond)
	old := Message{Path: "secret/old", At: time.Now().Add(-time.Hour)}
	// Manually enqueue an aged message.
	r.mu.Lock()
	r.queue = append(r.queue, old)
	r.mu.Unlock()

	fail.err = nil // allow delivery this time
	_ = r.Send(Message{Path: "secret/new", At: time.Now()})
	if r.QueueLen() != 0 {
		t.Fatalf("expected queue drained, got %d", r.QueueLen())
	}
}
