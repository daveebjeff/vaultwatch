package notify

import (
	"errors"
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

func TestNewDigestNotifier_NilInner(t *testing.T) {
	_, err := NewDigestNotifier(nil, time.Second, 10)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewDigestNotifier_ZeroWindow(t *testing.T) {
	_, err := NewDigestNotifier(&captureNotifier{}, 0, 10)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestDigestNotifier_FlushOnMaxSize(t *testing.T) {
	cap := &captureNotifier{}
	d, err := NewDigestNotifier(cap, 10*time.Second, 3)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		_ = d.Send(Message{Path: "secret/a", Status: StatusExpiringSoon, Detail: "soon"})
	}
	cap.mu.Lock()
	defer cap.mu.Unlock()
	if len(cap.msgs) != 1 {
		t.Fatalf("expected 1 digest message, got %d", len(cap.msgs))
	}
	if cap.msgs[0].Path != "digest" {
		t.Errorf("expected path 'digest', got %q", cap.msgs[0].Path)
	}
}

func TestDigestNotifier_ManualFlush(t *testing.T) {
	cap := &captureNotifier{}
	d, err := NewDigestNotifier(cap, 10*time.Second, 50)
	if err != nil {
		t.Fatal(err)
	}
	_ = d.Send(Message{Path: "secret/b", Status: StatusExpired, Detail: "expired"})
	if err := d.Flush(); err != nil {
		t.Fatal(err)
	}
	if len(cap.msgs) != 1 {
		t.Fatalf("expected 1 message after flush, got %d", len(cap.msgs))
	}
}

func TestDigestNotifier_FlushEmpty(t *testing.T) {
	cap := &captureNotifier{}
	d, _ := NewDigestNotifier(cap, time.Second, 10)
	if err := d.Flush(); err != nil {
		t.Errorf("flush on empty should not error: %v", err)
	}
	if len(cap.msgs) != 0 {
		t.Errorf("expected no messages, got %d", len(cap.msgs))
	}
}

func TestDigestNotifier_InnerError(t *testing.T) {
	cap := &captureNotifier{err: errors.New("send failed")}
	d, _ := NewDigestNotifier(cap, time.Second, 1)
	err := d.Send(Message{Path: "secret/c", Status: StatusExpired, Detail: "x"})
	if err == nil {
		t.Error("expected error from inner notifier")
	}
}

func TestDigestNotifier_WorstStatusPropagated(t *testing.T) {
	cap := &captureNotifier{}
	d, _ := NewDigestNotifier(cap, 10*time.Second, 50)
	_ = d.Send(Message{Path: "a", Status: StatusExpiringSoon, Detail: ""})
	_ = d.Send(Message{Path: "b", Status: StatusExpired, Detail: ""})
	_ = d.Flush()
	if cap.msgs[0].Status != StatusExpired {
		t.Errorf("expected StatusExpired, got %v", cap.msgs[0].Status)
	}
}
