package notify

import (
	"context"
	"testing"
	"time"
)

func TestNewSeenNotifier_NilInner(t *testing.T) {
	_, err := NewSeenNotifier(nil, time.Minute)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewSeenNotifier_ZeroWindow(t *testing.T) {
	_, err := NewSeenNotifier(NewNoopNotifier(), 0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewSeenNotifier_Valid(t *testing.T) {
	s, err := NewSeenNotifier(NewNoopNotifier(), time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestSeenNotifier_FirstSendForwarded(t *testing.T) {
	var count int
	inner := &mockNotifier{fn: func(_ context.Context, _ Message) error {
		count++
		return nil
	}}
	s, _ := NewSeenNotifier(inner, time.Hour)
	msg := Message{Path: "secret/a", Status: StatusExpiringSoon}

	if err := s.Send(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 send, got %d", count)
	}
}

func TestSeenNotifier_DuplicateSuppressed(t *testing.T) {
	var count int
	inner := &mockNotifier{fn: func(_ context.Context, _ Message) error {
		count++
		return nil
	}}
	s, _ := NewSeenNotifier(inner, time.Hour)
	msg := Message{Path: "secret/a", Status: StatusExpiringSoon}

	_ = s.Send(context.Background(), msg)
	_ = s.Send(context.Background(), msg)
	_ = s.Send(context.Background(), msg)

	if count != 1 {
		t.Fatalf("expected 1 send, got %d", count)
	}
}

func TestSeenNotifier_WindowExpiry_Reforwards(t *testing.T) {
	var count int
	inner := &mockNotifier{fn: func(_ context.Context, _ Message) error {
		count++
		return nil
	}}
	s, _ := NewSeenNotifier(inner, 10*time.Millisecond)
	msg := Message{Path: "secret/b", Status: StatusExpired}

	_ = s.Send(context.Background(), msg)
	time.Sleep(20 * time.Millisecond)
	_ = s.Send(context.Background(), msg)

	if count != 2 {
		t.Fatalf("expected 2 sends after window expiry, got %d", count)
	}
}

func TestSeenNotifier_Forget_AllowsResend(t *testing.T) {
	var count int
	inner := &mockNotifier{fn: func(_ context.Context, _ Message) error {
		count++
		return nil
	}}
	s, _ := NewSeenNotifier(inner, time.Hour)
	msg := Message{Path: "secret/c", Status: StatusExpiringSoon}

	_ = s.Send(context.Background(), msg)
	s.Forget(msg.Path)
	_ = s.Send(context.Background(), msg)

	if count != 2 {
		t.Fatalf("expected 2 sends after Forget, got %d", count)
	}
}

func TestSeenNotifier_Reset_AllowsAllResends(t *testing.T) {
	var count int
	inner := &mockNotifier{fn: func(_ context.Context, _ Message) error {
		count++
		return nil
	}}
	s, _ := NewSeenNotifier(inner, time.Hour)

	paths := []string{"secret/x", "secret/y", "secret/z"}
	for _, p := range paths {
		_ = s.Send(context.Background(), Message{Path: p, Status: StatusExpired})
	}
	s.Reset()
	for _, p := range paths {
		_ = s.Send(context.Background(), Message{Path: p, Status: StatusExpired})
	}

	if count != 6 {
		t.Fatalf("expected 6 sends after Reset, got %d", count)
	}
}

func TestSeenNotifier_IndependentPaths(t *testing.T) {
	var count int
	inner := &mockNotifier{fn: func(_ context.Context, _ Message) error {
		count++
		return nil
	}}
	s, _ := NewSeenNotifier(inner, time.Hour)

	_ = s.Send(context.Background(), Message{Path: "secret/p1", Status: StatusExpired})
	_ = s.Send(context.Background(), Message{Path: "secret/p2", Status: StatusExpired})
	_ = s.Send(context.Background(), Message{Path: "secret/p1", Status: StatusExpired})

	if count != 2 {
		t.Fatalf("expected 2 sends for independent paths, got %d", count)
	}
}
