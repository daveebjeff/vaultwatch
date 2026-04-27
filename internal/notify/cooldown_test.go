package notify

import (
	"context"
	"testing"
	"time"
)

func TestNewCooldownNotifier_NilInner(t *testing.T) {
	_, err := NewCooldownNotifier(nil, time.Minute)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewCooldownNotifier_ZeroDuration(t *testing.T) {
	n := NewNoopNotifier()
	_, err := NewCooldownNotifier(n, 0)
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
}

func TestNewCooldownNotifier_NegativeDuration(t *testing.T) {
	n := NewNoopNotifier()
	_, err := NewCooldownNotifier(n, -time.Second)
	if err == nil {
		t.Fatal("expected error for negative duration")
	}
}

func TestNewCooldownNotifier_Valid(t *testing.T) {
	n := NewNoopNotifier()
	c, err := NewCooldownNotifier(n, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestCooldownNotifier_FirstSendForwarded(t *testing.T) {
	mock := &mockNotifier{}
	c, _ := NewCooldownNotifier(mock, time.Minute)

	msg := Message{Path: "secret/foo", Status: StatusExpiringSoon}
	if err := c.Send(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.calls != 1 {
		t.Fatalf("expected 1 call, got %d", mock.calls)
	}
}

func TestCooldownNotifier_SecondSendSuppressed(t *testing.T) {
	mock := &mockNotifier{}
	c, _ := NewCooldownNotifier(mock, time.Minute)

	msg := Message{Path: "secret/foo", Status: StatusExpiringSoon}
	_ = c.Send(context.Background(), msg)
	_ = c.Send(context.Background(), msg)

	if mock.calls != 1 {
		t.Fatalf("expected 1 call after suppression, got %d", mock.calls)
	}
}

func TestCooldownNotifier_DifferentPathsIndependent(t *testing.T) {
	mock := &mockNotifier{}
	c, _ := NewCooldownNotifier(mock, time.Minute)

	_ = c.Send(context.Background(), Message{Path: "secret/a", Status: StatusExpiringSoon})
	_ = c.Send(context.Background(), Message{Path: "secret/b", Status: StatusExpiringSoon})

	if mock.calls != 2 {
		t.Fatalf("expected 2 calls for different paths, got %d", mock.calls)
	}
}

func TestCooldownNotifier_ForwardsAfterCooldown(t *testing.T) {
	mock := &mockNotifier{}
	c, _ := NewCooldownNotifier(mock, 10*time.Millisecond)

	msg := Message{Path: "secret/foo", Status: StatusExpiringSoon}
	_ = c.Send(context.Background(), msg)

	time.Sleep(20 * time.Millisecond)
	_ = c.Send(context.Background(), msg)

	if mock.calls != 2 {
		t.Fatalf("expected 2 calls after cooldown elapsed, got %d", mock.calls)
	}
}
