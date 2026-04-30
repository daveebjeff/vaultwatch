package notify

import (
	"errors"
	"testing"
	"time"
)

func sieveMsg(path string, status Status) Message {
	return Message{
		Path:      path,
		Status:    status,
		ExpiresAt: time.Now().Add(time.Hour),
	}
}

func TestNewSieveNotifier_Empty(t *testing.T) {
	s := NewSieveNotifier()
	if s == nil {
		t.Fatal("expected non-nil SieveNotifier")
	}
}

func TestSieveNotifier_Add_NilPredicate(t *testing.T) {
	s := NewSieveNotifier()
	err := s.Add(nil, NewNoopNotifier())
	if err == nil {
		t.Fatal("expected error for nil predicate")
	}
}

func TestSieveNotifier_Add_NilNotifier(t *testing.T) {
	s := NewSieveNotifier()
	err := s.Add(func(Message) bool { return true }, nil)
	if err == nil {
		t.Fatal("expected error for nil notifier")
	}
}

func TestSieveNotifier_NoRoutes_ReturnsNil(t *testing.T) {
	s := NewSieveNotifier()
	if err := s.Send(sieveMsg("secret/a", StatusExpired)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSieveNotifier_RoutesToFirstMatch(t *testing.T) {
	s := NewSieveNotifier()

	var gotA, gotB bool
	nA := &mockNotifier{fn: func(Message) error { gotA = true; return nil }}
	nB := &mockNotifier{fn: func(Message) error { gotB = true; return nil }}

	_ = s.Add(func(m Message) bool { return m.Status == StatusExpired }, nA)
	_ = s.Add(func(m Message) bool { return true }, nB)

	if err := s.Send(sieveMsg("secret/x", StatusExpired)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !gotA {
		t.Error("expected first notifier to be called")
	}
	if gotB {
		t.Error("expected second notifier NOT to be called")
	}
}

func TestSieveNotifier_FallsThroughToSecond(t *testing.T) {
	s := NewSieveNotifier()

	var gotB bool
	nA := &mockNotifier{fn: func(Message) error { return nil }}
	nB := &mockNotifier{fn: func(Message) error { gotB = true; return nil }}

	_ = s.Add(func(m Message) bool { return m.Status == StatusExpired }, nA)
	_ = s.Add(func(m Message) bool { return true }, nB)

	if err := s.Send(sieveMsg("secret/y", StatusExpiringSoon)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !gotB {
		t.Error("expected second notifier to be called")
	}
}

func TestSieveNotifier_PropagatesError(t *testing.T) {
	s := NewSieveNotifier()
	sentinel := errors.New("send failed")
	n := &mockNotifier{fn: func(Message) error { return sentinel }}
	_ = s.Add(func(Message) bool { return true }, n)

	if err := s.Send(sieveMsg("secret/z", StatusExpired)); !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}
