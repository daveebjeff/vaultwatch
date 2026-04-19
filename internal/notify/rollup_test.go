package notify

import (
	"fmt"
	"strings"
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
	if c.err != nil {
		return c.err
	}
	c.msgs = append(c.msgs, m)
	return nil
}

func TestNewRollupNotifier_NilInner(t *testing.T) {
	_, err := NewRollupNotifier(nil, time.Second, 10)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewRollupNotifier_ZeroWindow(t *testing.T) {
	c := &captureNotifier{}
	_, err := NewRollupNotifier(c, 0, 10)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestRollupNotifier_FlushOnMaxSize(t *testing.T) {
	c := &captureNotifier{}
	r, err := NewRollupNotifier(c, 10*time.Second, 3)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		_ = r.Send(Message{SecretPath: fmt.Sprintf("secret/%d", i), Status: StatusExpiringSoon})
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.msgs) != 1 {
		t.Fatalf("expected 1 rollup message, got %d", len(c.msgs))
	}
	if !strings.Contains(c.msgs[0].Summary, "3 secret(s)") {
		t.Errorf("unexpected summary: %s", c.msgs[0].Summary)
	}
}

func TestRollupNotifier_ManualFlush(t *testing.T) {
	c := &captureNotifier{}
	r, _ := NewRollupNotifier(c, 10*time.Second, 50)
	_ = r.Send(Message{SecretPath: "secret/a", Status: StatusExpired})
	_ = r.Send(Message{SecretPath: "secret/b", Status: StatusExpiringSoon})
	if err := r.Flush(); err != nil {
		t.Fatal(err)
	}
	if len(c.msgs) != 1 {
		t.Fatalf("expected 1 rollup, got %d", len(c.msgs))
	}
	if c.msgs[0].Status != StatusExpired {
		t.Errorf("expected worst status Expired, got %v", c.msgs[0].Status)
	}
}

func TestRollupNotifier_FlushEmpty(t *testing.T) {
	c := &captureNotifier{}
	r, _ := NewRollupNotifier(c, time.Second, 10)
	if err := r.Flush(); err != nil {
		t.Errorf("flush on empty should not error: %v", err)
	}
	if len(c.msgs) != 0 {
		t.Errorf("expected no messages sent")
	}
}

func TestRollupNotifier_WindowFlush(t *testing.T) {
	c := &captureNotifier{}
	r, _ := NewRollupNotifier(c, 50*time.Millisecond, 100)
	_ = r.Send(Message{SecretPath: "secret/x", Status: StatusExpiringSoon})
	time.Sleep(150 * time.Millisecond)
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.msgs) != 1 {
		t.Fatalf("expected 1 message after window, got %d", len(c.msgs))
	}
}
