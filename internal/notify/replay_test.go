package notify

import (
	"errors"
	"testing"
	"time"
)

func replayMsg(path string, offset time.Duration) Message {
	return Message{
		Path:      path,
		Status:    StatusExpiringSoon,
		ExpiresAt: time.Now().Add(offset),
	}
}

func TestNewReplayNotifier_NilInner(t *testing.T) {
	_, err := NewReplayNotifier(nil, 10, time.Minute)
	if !errors.Is(err, ErrNilInner) {
		t.Fatalf("expected ErrNilInner, got %v", err)
	}
}

func TestNewReplayNotifier_ZeroMax(t *testing.T) {
	n := NewNoopNotifier()
	_, err := NewReplayNotifier(n, 0, time.Minute)
	if err == nil {
		t.Fatal("expected error for zero maxItems")
	}
}

func TestNewReplayNotifier_ZeroMaxAge(t *testing.T) {
	n := NewNoopNotifier()
	_, err := NewReplayNotifier(n, 5, 0)
	if err == nil {
		t.Fatal("expected error for zero maxAge")
	}
}

func TestNewReplayNotifier_Valid(t *testing.T) {
	n := NewNoopNotifier()
	r, err := NewReplayNotifier(n, 5, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestReplayNotifier_Send_ForwardsToInner(t *testing.T) {
	var got []Message
	collector := &mockNotifier{fn: func(m Message) error {
		got = append(got, m)
		return nil
	}}
	r, _ := NewReplayNotifier(collector, 10, time.Hour)
	msg := replayMsg("secret/foo", time.Hour)
	if err := r.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Path != "secret/foo" {
		t.Fatalf("inner not called correctly, got %v", got)
	}
}

func TestReplayNotifier_Replay_ForwardsRetained(t *testing.T) {
	r, _ := NewReplayNotifier(NewNoopNotifier(), 10, time.Hour)
	r.Send(replayMsg("secret/a", time.Hour))
	r.Send(replayMsg("secret/b", time.Hour))

	var replayed []Message
	dst := &mockNotifier{fn: func(m Message) error {
		replayed = append(replayed, m)
		return nil
	}}
	if err := r.Replay(dst); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(replayed) != 2 {
		t.Fatalf("expected 2 replayed messages, got %d", len(replayed))
	}
}

func TestReplayNotifier_MaxItemsRespected(t *testing.T) {
	r, _ := NewReplayNotifier(NewNoopNotifier(), 3, time.Hour)
	for i := 0; i < 5; i++ {
		r.Send(replayMsg("secret/x", time.Hour))
	}
	if r.Len() != 3 {
		t.Fatalf("expected 3 retained messages, got %d", r.Len())
	}
}

func TestReplayNotifier_Replay_ErrorPropagates(t *testing.T) {
	r, _ := NewReplayNotifier(NewNoopNotifier(), 10, time.Hour)
	r.Send(replayMsg("secret/a", time.Hour))

	failDst := &mockNotifier{fn: func(m Message) error {
		return errors.New("dst failure")
	}}
	if err := r.Replay(failDst); err == nil {
		t.Fatal("expected error from failing dst")
	}
}
