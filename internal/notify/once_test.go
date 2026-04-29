package notify

import (
	"context"
	"testing"
	"time"
)

func onceMsg(path string) Message {
	return Message{
		Path:      path,
		Status:    StatusExpiringSoon,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
}

func TestNewOnceNotifier_NilInner(t *testing.T) {
	_, err := NewOnceNotifier(nil)
	if err == nil {
		t.Fatal("expected error for nil inner, got nil")
	}
}

func TestNewOnceNotifier_Valid(t *testing.T) {
	noop, _ := NewNoopNotifier()
	n, err := NewOnceNotifier(noop)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil OnceNotifier")
	}
}

func TestOnceNotifier_FirstSendForwarded(t *testing.T) {
	mock := &mockNotifier{}
	n, _ := NewOnceNotifier(mock)

	if err := n.Send(context.Background(), onceMsg("secret/a")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.calls != 1 {
		t.Fatalf("expected 1 call, got %d", mock.calls)
	}
}

func TestOnceNotifier_SecondSendSuppressed(t *testing.T) {
	mock := &mockNotifier{}
	n, _ := NewOnceNotifier(mock)

	msg := onceMsg("secret/a")
	_ = n.Send(context.Background(), msg)
	_ = n.Send(context.Background(), msg)

	if mock.calls != 1 {
		t.Fatalf("expected 1 call after duplicate, got %d", mock.calls)
	}
}

func TestOnceNotifier_DifferentPathsForwarded(t *testing.T) {
	mock := &mockNotifier{}
	n, _ := NewOnceNotifier(mock)

	_ = n.Send(context.Background(), onceMsg("secret/a"))
	_ = n.Send(context.Background(), onceMsg("secret/b"))

	if mock.calls != 2 {
		t.Fatalf("expected 2 calls for distinct paths, got %d", mock.calls)
	}
}

func TestOnceNotifier_ResetAllowsResend(t *testing.T) {
	mock := &mockNotifier{}
	n, _ := NewOnceNotifier(mock)

	_ = n.Send(context.Background(), onceMsg("secret/a"))
	n.Reset()
	_ = n.Send(context.Background(), onceMsg("secret/a"))

	if mock.calls != 2 {
		t.Fatalf("expected 2 calls after reset, got %d", mock.calls)
	}
}

func TestOnceNotifier_Seen(t *testing.T) {
	mock := &mockNotifier{}
	n, _ := NewOnceNotifier(mock)

	if n.Seen("secret/a") {
		t.Fatal("expected path not seen before first send")
	}
	_ = n.Send(context.Background(), onceMsg("secret/a"))
	if !n.Seen("secret/a") {
		t.Fatal("expected path seen after first send")
	}
}
