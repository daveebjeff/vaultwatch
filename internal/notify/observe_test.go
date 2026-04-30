package notify

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewObserveNotifier_NilInner(t *testing.T) {
	_, err := NewObserveNotifier(nil)
	if err == nil {
		t.Fatal("expected error for nil inner, got nil")
	}
}

func TestNewObserveNotifier_Valid(t *testing.T) {
	n, err := NewObserveNotifier(NewNoopNotifier())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil ObserveNotifier")
	}
}

func TestObserveNotifier_CountsSuccess(t *testing.T) {
	n, _ := NewObserveNotifier(NewNoopNotifier())
	msg := Message{Path: "secret/foo", Status: StatusExpiringSoon}

	for i := 0; i < 3; i++ {
		if err := n.Send(context.Background(), msg); err != nil {
			t.Fatalf("unexpected send error: %v", err)
		}
	}

	s := n.Stats()
	if s.Total != 3 {
		t.Errorf("expected Total=3, got %d", s.Total)
	}
	if s.Errors != 0 {
		t.Errorf("expected Errors=0, got %d", s.Errors)
	}
}

func TestObserveNotifier_CountsErrors(t *testing.T) {
	sentinel := errors.New("boom")
	failing := &mockNotifier{err: sentinel}
	n, _ := NewObserveNotifier(failing)
	msg := Message{Path: "secret/bar", Status: StatusExpired}

	_ = n.Send(context.Background(), msg)
	_ = n.Send(context.Background(), msg)

	s := n.Stats()
	if s.Total != 2 {
		t.Errorf("expected Total=2, got %d", s.Total)
	}
	if s.Errors != 2 {
		t.Errorf("expected Errors=2, got %d", s.Errors)
	}
}

func TestObserveNotifier_AvgLatencyNonZero(t *testing.T) {
	slow := &delayNotifier{delay: 5 * time.Millisecond}
	n, _ := NewObserveNotifier(slow)
	msg := Message{Path: "secret/baz", Status: StatusExpiringSoon}

	_ = n.Send(context.Background(), msg)

	s := n.Stats()
	if s.AvgLatency <= 0 {
		t.Errorf("expected positive avg latency, got %v", s.AvgLatency)
	}
}

func TestObserveNotifier_Reset(t *testing.T) {
	n, _ := NewObserveNotifier(NewNoopNotifier())
	msg := Message{Path: "secret/foo", Status: StatusExpiringSoon}
	_ = n.Send(context.Background(), msg)

	n.Reset()
	s := n.Stats()
	if s.Total != 0 || s.Errors != 0 || s.AvgLatency != 0 {
		t.Errorf("expected zeroed stats after Reset, got %+v", s)
	}
}

// delayNotifier is a test helper that sleeps before returning.
type delayNotifier struct{ delay time.Duration }

func (d *delayNotifier) Send(_ context.Context, _ Message) error {
	time.Sleep(d.delay)
	return nil
}
