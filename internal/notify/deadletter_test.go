package notify

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewDeadLetterNotifier_NilInner(t *testing.T) {
	_, err := NewDeadLetterNotifier(nil, 10)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewDeadLetterNotifier_Valid(t *testing.T) {
	n, err := NewDeadLetterNotifier(NewNoopNotifier(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestDeadLetterNotifier_SuccessNotCaptured(t *testing.T) {
	n, _ := NewDeadLetterNotifier(NewNoopNotifier(), 10)
	msg := Message{Path: "secret/a", Status: StatusExpiringSoon, ExpiresAt: time.Now().Add(time.Hour)}
	_ = n.Send(context.Background(), msg)
	if len(n.Failed()) != 0 {
		t.Fatal("expected no dead letters on success")
	}
}

func TestDeadLetterNotifier_FailureCaptured(t *testing.T) {
	sentinel := errors.New("send failed")
	inner := &mockFailNotifier{err: sentinel}
	n, _ := NewDeadLetterNotifier(inner, 10)
	msg := Message{Path: "secret/b", Status: StatusExpired, ExpiresAt: time.Now()}
	err := n.Send(context.Background(), msg)
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	failed := n.Failed()
	if len(failed) != 1 {
		t.Fatalf("expected 1 dead letter, got %d", len(failed))
	}
	if failed[0].Message.Path != "secret/b" {
		t.Errorf("unexpected path: %s", failed[0].Message.Path)
	}
}

func TestDeadLetterNotifier_MaxSizeRespected(t *testing.T) {
	inner := &mockFailNotifier{err: errors.New("fail")}
	n, _ := NewDeadLetterNotifier(inner, 3)
	msg := Message{Path: "secret/x", Status: StatusExpired}
	for i := 0; i < 10; i++ {
		_ = n.Send(context.Background(), msg)
	}
	if len(n.Failed()) != 3 {
		t.Fatalf("expected 3 dead letters, got %d", len(n.Failed()))
	}
}

func TestDeadLetterNotifier_Drain(t *testing.T) {
	inner := &mockFailNotifier{err: errors.New("fail")}
	n, _ := NewDeadLetterNotifier(inner, 10)
	msg := Message{Path: "secret/y", Status: StatusExpired}
	_ = n.Send(context.Background(), msg)
	drained := n.Drain()
	if len(drained) != 1 {
		t.Fatalf("expected 1 drained entry, got %d", len(drained))
	}
	if len(n.Failed()) != 0 {
		t.Fatal("expected empty queue after drain")
	}
}
