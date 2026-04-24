package notify

import (
	"sync"
	"testing"
	"time"
)

func TestNewDebounceNotifier_NilInner(t *testing.T) {
	_, err := NewDebounceNotifier(nil, 10*time.Millisecond)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewDebounceNotifier_ZeroWait(t *testing.T) {
	base, _ := NewNoopNotifier()
	_, err := NewDebounceNotifier(base, 0)
	if err == nil {
		t.Fatal("expected error for zero wait")
	}
}

func TestNewDebounceNotifier_Valid(t *testing.T) {
	base, _ := NewNoopNotifier()
	d, err := NewDebounceNotifier(base, 10*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestDebounceNotifier_OnlyLastMessageForwarded(t *testing.T) {
	var mu sync.Mutex
	var received []Message

	collector := &mockNotifier{
		sendFn: func(m Message) error {
			mu.Lock()
			received = append(received, m)
			mu.Unlock()
			return nil
		},
	}

	wait := 40 * time.Millisecond
	d, _ := NewDebounceNotifier(collector, wait)

	for i := 0; i < 5; i++ {
		_ = d.Send(Message{Path: "secret/foo", Body: string(rune('a' + i))})
		time.Sleep(5 * time.Millisecond)
	}

	time.Sleep(wait + 20*time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(received) != 1 {
		t.Fatalf("expected 1 forwarded message, got %d", len(received))
	}
	if received[0].Body != "e" {
		t.Errorf("expected last body 'e', got %q", received[0].Body)
	}
}

func TestDebounceNotifier_IndependentPathTimers(t *testing.T) {
	var mu sync.Mutex
	var received []Message

	collector := &mockNotifier{
		sendFn: func(m Message) error {
			mu.Lock()
			received = append(received, m)
			mu.Unlock()
			return nil
		},
	}

	wait := 30 * time.Millisecond
	d, _ := NewDebounceNotifier(collector, wait)

	_ = d.Send(Message{Path: "secret/a", Body: "msgA"})
	_ = d.Send(Message{Path: "secret/b", Body: "msgB"})

	time.Sleep(wait + 20*time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(received) != 2 {
		t.Fatalf("expected 2 messages (one per path), got %d", len(received))
	}
}
