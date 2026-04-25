package notify

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewTapNotifier_NilInner(t *testing.T) {
	_, err := NewTapNotifier(nil, func(Message) {})
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewTapNotifier_NilFn(t *testing.T) {
	_, err := NewTapNotifier(NewNoopNotifier(), nil)
	if err == nil {
		t.Fatal("expected error for nil tap function")
	}
}

func TestNewTapNotifier_Valid(t *testing.T) {
	n, err := NewTapNotifier(NewNoopNotifier(), func(Message) {})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestTapNotifier_TapReceivesMessage(t *testing.T) {
	var got Message
	tapped := make(chan struct{}, 1)

	n, _ := NewTapNotifier(NewNoopNotifier(), func(m Message) {
		got = m
		tapped <- struct{}{}
	})

	want := Message{
		Path:      "secret/db",
		Status:    StatusExpiringSoon,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	if err := n.Send(context.Background(), want); err != nil {
		t.Fatalf("Send returned error: %v", err)
	}

	select {
	case <-tapped:
	case <-time.After(time.Second):
		t.Fatal("tap function was not called")
	}

	if got.Path != want.Path {
		t.Errorf("got path %q, want %q", got.Path, want.Path)
	}
	if got.Status != want.Status {
		t.Errorf("got status %v, want %v", got.Status, want.Status)
	}
}

func TestTapNotifier_InnerErrorReturned(t *testing.T) {
	sentinel := errors.New("inner failure")
	failing := &mockNotifier{err: sentinel}

	n, _ := NewTapNotifier(failing, func(Message) {})
	err := n.Send(context.Background(), Message{Path: "x"})
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestTapNotifier_PanicInTapSuppressed(t *testing.T) {
	n, _ := NewTapNotifier(NewNoopNotifier(), func(Message) {
		panic("tap exploded")
	})

	// Should not panic; inner Send should still succeed.
	if err := n.Send(context.Background(), Message{Path: "y"}); err != nil {
		t.Fatalf("unexpected error after tap panic: %v", err)
	}
}

func TestTapNotifier_SetTapFn_Replaces(t *testing.T) {
	calls := make([]string, 0, 2)

	n, _ := NewTapNotifier(NewNoopNotifier(), func(m Message) {
		calls = append(calls, "first")
	})
	_ = n.Send(context.Background(), Message{Path: "a"})

	n.SetTapFn(func(m Message) {
		calls = append(calls, "second")
	})
	_ = n.Send(context.Background(), Message{Path: "b"})

	if len(calls) != 2 || calls[0] != "first" || calls[1] != "second" {
		t.Errorf("unexpected call sequence: %v", calls)
	}
}
