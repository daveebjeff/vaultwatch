package notify

import (
	"context"
	"errors"
	"testing"
)

func TestNewPassthroughNotifier_NilInner(t *testing.T) {
	_, err := NewPassthroughNotifier(nil)
	if err == nil {
		t.Fatal("expected error for nil inner, got nil")
	}
}

func TestNewPassthroughNotifier_Valid(t *testing.T) {
	p, err := NewPassthroughNotifier(NewNoopNotifier())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil PassthroughNotifier")
	}
}

func TestPassthroughNotifier_CountsSuccess(t *testing.T) {
	p, _ := NewPassthroughNotifier(NewNoopNotifier())
	msg := Message{Path: "secret/foo", Status: StatusExpiringSoon}

	for i := 0; i < 3; i++ {
		if err := p.Send(context.Background(), msg); err != nil {
			t.Fatalf("unexpected error on send %d: %v", i, err)
		}
	}

	if p.Seen() != 3 {
		t.Errorf("Seen: want 3, got %d", p.Seen())
	}
	if p.Sent() != 3 {
		t.Errorf("Sent: want 3, got %d", p.Sent())
	}
	if p.Errors() != 0 {
		t.Errorf("Errors: want 0, got %d", p.Errors())
	}
}

func TestPassthroughNotifier_CountsErrors(t *testing.T) {
	failing := &mockNotifier{err: errors.New("boom")}
	p, _ := NewPassthroughNotifier(failing)
	msg := Message{Path: "secret/bar", Status: StatusExpired}

	_ = p.Send(context.Background(), msg)
	_ = p.Send(context.Background(), msg)

	if p.Seen() != 2 {
		t.Errorf("Seen: want 2, got %d", p.Seen())
	}
	if p.Sent() != 0 {
		t.Errorf("Sent: want 0, got %d", p.Sent())
	}
	if p.Errors() != 2 {
		t.Errorf("Errors: want 2, got %d", p.Errors())
	}
}

func TestPassthroughNotifier_Reset(t *testing.T) {
	p, _ := NewPassthroughNotifier(NewNoopNotifier())
	msg := Message{Path: "secret/baz", Status: StatusOK}

	_ = p.Send(context.Background(), msg)
	p.Reset()

	if p.Seen() != 0 || p.Sent() != 0 || p.Errors() != 0 {
		t.Errorf("Reset did not zero counters: seen=%d sent=%d errors=%d",
			p.Seen(), p.Sent(), p.Errors())
	}
}

func TestPassthroughNotifier_ForwardsMessage(t *testing.T) {
	var received []Message
	capture := &mockNotifier{
		sendFn: func(_ context.Context, m Message) error {
			received = append(received, m)
			return nil
		},
	}
	p, _ := NewPassthroughNotifier(capture)
	want := Message{Path: "secret/capture", Status: StatusExpired, Body: "expired"}

	if err := p.Send(context.Background(), want); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received) != 1 {
		t.Fatalf("expected 1 message forwarded, got %d", len(received))
	}
	if received[0] != want {
		t.Errorf("forwarded message mismatch: got %+v, want %+v", received[0], want)
	}
}
