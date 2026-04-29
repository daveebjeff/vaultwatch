package notify

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewStaggerNotifier_NoNotifiers(t *testing.T) {
	_, err := NewStaggerNotifier(10*time.Millisecond)
	if err == nil {
		t.Fatal("expected error for no notifiers")
	}
}

func TestNewStaggerNotifier_NilNotifier(t *testing.T) {
	_, err := NewStaggerNotifier(10*time.Millisecond, NewNoopNotifier(), nil)
	if err == nil {
		t.Fatal("expected error for nil notifier")
	}
}

func TestNewStaggerNotifier_ZeroDelay(t *testing.T) {
	_, err := NewStaggerNotifier(0, NewNoopNotifier())
	if err == nil {
		t.Fatal("expected error for zero delay")
	}
}

func TestNewStaggerNotifier_Valid(t *testing.T) {
	s, err := NewStaggerNotifier(5*time.Millisecond, NewNoopNotifier())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil StaggerNotifier")
	}
}

func TestStaggerNotifier_AllCalled(t *testing.T) {
	var calls []int
	make := func(id int) Notifier {
		return &mockNotifier{fn: func(_ context.Context, _ Message) error {
			calls = append(calls, id)
			return nil
		}}
	}
	s, _ := NewStaggerNotifier(5*time.Millisecond, make(1), make(2), make(3))
	if err := s.Send(context.Background(), Message{Path: "x"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(calls) != 3 {
		t.Fatalf("expected 3 calls, got %d", len(calls))
	}
	for i, v := range calls {
		if v != i+1 {
			t.Errorf("expected call %d got %d", i+1, v)
		}
	}
}

func TestStaggerNotifier_ReturnsFirstError(t *testing.T) {
	sentinel := errors.New("boom")
	noop := NewNoopNotifier()
	failing := &mockNotifier{fn: func(_ context.Context, _ Message) error { return sentinel }}
	s, _ := NewStaggerNotifier(5*time.Millisecond, failing, noop)
	err := s.Send(context.Background(), Message{Path: "p"})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestStaggerNotifier_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	n1 := &mockNotifier{fn: func(_ context.Context, _ Message) error { return nil }}
	s, _ := NewStaggerNotifier(50*time.Millisecond, n1, n1)
	err := s.Send(ctx, Message{Path: "p"})
	if err == nil {
		t.Fatal("expected context error")
	}
}

func TestStaggerNotifier_Add(t *testing.T) {
	s, _ := NewStaggerNotifier(5*time.Millisecond, NewNoopNotifier())
	if err := s.Add(NewNoopNotifier()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.inners) != 2 {
		t.Fatalf("expected 2 inners, got %d", len(s.inners))
	}
}

func TestStaggerNotifier_Add_Nil(t *testing.T) {
	s, _ := NewStaggerNotifier(5*time.Millisecond, NewNoopNotifier())
	if err := s.Add(nil); err == nil {
		t.Fatal("expected error adding nil notifier")
	}
}
