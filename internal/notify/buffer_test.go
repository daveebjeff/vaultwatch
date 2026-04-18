package notify

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type captureNotifier struct {
	mu   sync.Mutex
	msgs []Message
	err  error
}

func (c *captureNotifier) Send(m Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.msgs = append(c.msgs, m)
	return c.err
}

func (c *captureNotifier) count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.msgs)
}

func TestNewBufferNotifier_NilInner(t *testing.T) {
	_, err := NewBufferNotifier(nil, time.Second, 0)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewBufferNotifier_ZeroWindow(t *testing.T) {
	cap := &captureNotifier{}
	_, err := NewBufferNotifier(cap, 0, 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestBufferNotifier_FlushOnMaxSize(t *testing.T) {
	cap := &captureNotifier{}
	b, err := NewBufferNotifier(cap, 10*time.Second, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	msg := Message{Path: "secret/foo", Status: StatusExpired}
	for i := 0; i < 2; i++ {
		_ = b.Send(msg)
	}
	if cap.count() != 0 {
		t.Fatal("should not have flushed yet")
	}
	_ = b.Send(msg)
	if cap.count() != 3 {
		t.Fatalf("expected 3 sent, got %d", cap.count())
	}
}

func TestBufferNotifier_ManualFlush(t *testing.T) {
	cap := &captureNotifier{}
	b, _ := NewBufferNotifier(cap, 10*time.Second, 0)
	_ = b.Send(Message{Path: "secret/a", Status: StatusExpiringSoon})
	_ = b.Send(Message{Path: "secret/b", Status: StatusExpired})
	if err := b.Flush(); err != nil {
		t.Fatalf("flush error: %v", err)
	}
	if cap.count() != 2 {
		t.Fatalf("expected 2 sent, got %d", cap.count())
	}
}

func TestBufferNotifier_TimerFlush(t *testing.T) {
	cap := &captureNotifier{}
	b, _ := NewBufferNotifier(cap, 50*time.Millisecond, 0)
	_ = b.Send(Message{Path: "secret/timer", Status: StatusExpired})
	time.Sleep(150 * time.Millisecond)
	if cap.count() != 1 {
		t.Fatalf("expected 1 sent after timer, got %d", cap.count())
	}
}

func TestBufferNotifier_FlushReturnsError(t *testing.T) {
	cap := &captureNotifier{err: fmt.Errorf("send failed")}
	b, _ := NewBufferNotifier(cap, 10*time.Second, 0)
	_ = b.Send(Message{Path: "secret/err", Status: StatusExpired})
	if err := b.Flush(); err == nil {
		t.Fatal("expected error from flush")
	}
}
