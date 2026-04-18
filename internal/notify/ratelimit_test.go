package notify

import (
	"errors"
	"testing"
	"time"
)

func TestNewRateLimitNotifier_NilNotifier(t *testing.T) {
	_, err := NewRateLimitNotifier(nil, time.Minute)
	if err == nil {
		t.Fatal("expected error for nil notifier")
	}
}

func TestNewRateLimitNotifier_ZeroCooldown(t *testing.T) {
	_, err := NewRateLimitNotifier(&countingNotifier{}, 0)
	if err == nil {
		t.Fatal("expected error for zero cooldown")
	}
}

func TestRateLimitNotifier_FirstSendAllowed(t *testing.T) {
	n := &countingNotifier{}
	r, _ := NewRateLimitNotifier(n, time.Minute)
	if err := r.Send(exampleMsg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.calls != 1 {
		t.Errorf("expected 1 call, got %d", n.calls)
	}
}

func TestRateLimitNotifier_DuplicateSuppressed(t *testing.T) {
	n := &countingNotifier{}
	r, _ := NewRateLimitNotifier(n, time.Minute)
	_ = r.Send(exampleMsg)
	_ = r.Send(exampleMsg)
	if n.calls != 1 {
		t.Errorf("expected 1 call, got %d", n.calls)
	}
}

func TestRateLimitNotifier_DifferentPathsAllowed(t *testing.T) {
	n := &countingNotifier{}
	r, _ := NewRateLimitNotifier(n, time.Minute)
	msg1 := Message{Path: "secret/a"}
	msg2 := Message{Path: "secret/b"}
	_ = r.Send(msg1)
	_ = r.Send(msg2)
	if n.calls != 2 {
		t.Errorf("expected 2 calls, got %d", n.calls)
	}
}

func TestRateLimitNotifier_CooldownExpiry(t *testing.T) {
	n := &countingNotifier{}
	r, _ := NewRateLimitNotifier(n, 10*time.Millisecond)
	_ = r.Send(exampleMsg)
	time.Sleep(20 * time.Millisecond)
	_ = r.Send(exampleMsg)
	if n.calls != 2 {
		t.Errorf("expected 2 calls after cooldown, got %d", n.calls)
	}
}

func TestRateLimitNotifier_PropagatesError(t *testing.T) {
	n := &countingNotifier{failUntil: 1, err: errors.New("send failed")}
	r, _ := NewRateLimitNotifier(n, time.Minute)
	if err := r.Send(exampleMsg); err == nil {
		t.Fatal("expected error to propagate")
	}
}
