package notify

import (
	"errors"
	"testing"
)

func TestNewSamplingNotifier_NilInner(t *testing.T) {
	_, err := NewSamplingNotifier(nil, 0.5)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewSamplingNotifier_InvalidRate(t *testing.T) {
	n := NewNoopNotifier()
	for _, rate := range []float64{-0.1, 1.1, 2.0} {
		_, err := NewSamplingNotifier(n, rate)
		if err == nil {
			t.Fatalf("expected error for rate %v", rate)
		}
	}
}

func TestNewSamplingNotifier_Valid(t *testing.T) {
	n := NewNoopNotifier()
	s, err := NewSamplingNotifier(n, 0.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestSamplingNotifier_RateZero_NeverForwards(t *testing.T) {
	var called int
	fn := notifierFunc(func(_ Message) error {
		called++
		return nil
	})
	s, _ := NewSamplingNotifier(fn, 0.0)
	msg := Message{Path: "secret/test", Status: StatusExpired}
	for i := 0; i < 100; i++ {
		_ = s.Send(msg)
	}
	if called != 0 {
		t.Fatalf("expected 0 forwards, got %d", called)
	}
}

func TestSamplingNotifier_RateOne_AlwaysForwards(t *testing.T) {
	var called int
	fn := notifierFunc(func(_ Message) error {
		called++
		return nil
	})
	s, _ := NewSamplingNotifier(fn, 1.0)
	msg := Message{Path: "secret/test", Status: StatusExpired}
	for i := 0; i < 20; i++ {
		_ = s.Send(msg)
	}
	if called != 20 {
		t.Fatalf("expected 20 forwards, got %d", called)
	}
}

func TestSamplingNotifier_PropagatesError(t *testing.T) {
	sentinel := errors.New("send failed")
	fn := notifierFunc(func(_ Message) error { return sentinel })
	s, _ := NewSamplingNotifier(fn, 1.0)
	err := s.Send(Message{Path: "secret/x", Status: StatusExpiringSoon})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

// notifierFunc is a helper adapter already used in other test files.
type notifierFunc func(Message) error

func (f notifierFunc) Send(m Message) error { return f(m) }
