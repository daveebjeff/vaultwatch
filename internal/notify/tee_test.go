package notify

import (
	"context"
	"errors"
	"testing"
	"time"
)

func teeMsg() Message {
	return Message{
		Path:      "secret/tee",
		Status:    StatusExpiringSoon,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
}

func TestNewTeeNotifier_NilA(t *testing.T) {
	_, err := NewTeeNotifier(nil, &mockNotifier{})
	if err == nil {
		t.Fatal("expected error for nil first notifier")
	}
}

func TestNewTeeNotifier_NilB(t *testing.T) {
	_, err := NewTeeNotifier(&mockNotifier{}, nil)
	if err == nil {
		t.Fatal("expected error for nil second notifier")
	}
}

func TestNewTeeNotifier_Valid(t *testing.T) {
	tee, err := NewTeeNotifier(&mockNotifier{}, &mockNotifier{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tee == nil {
		t.Fatal("expected non-nil TeeNotifier")
	}
}

func TestTeeNotifier_BothCalled(t *testing.T) {
	a := &mockNotifier{}
	b := &mockNotifier{}
	tee, _ := NewTeeNotifier(a, b)

	if err := tee.Send(context.Background(), teeMsg()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.calls != 1 {
		t.Errorf("expected a to be called once, got %d", a.calls)
	}
	if b.calls != 1 {
		t.Errorf("expected b to be called once, got %d", b.calls)
	}
}

func TestTeeNotifier_BothCalledEvenIfAFails(t *testing.T) {
	a := &mockNotifier{err: errors.New("a failed")}
	b := &mockNotifier{}
	tee, _ := NewTeeNotifier(a, b)

	err := tee.Send(context.Background(), teeMsg())
	if err == nil {
		t.Fatal("expected error from failing notifier")
	}
	if b.calls != 1 {
		t.Errorf("b should still be called when a fails, got %d calls", b.calls)
	}
}

func TestTeeNotifier_BothErrorsJoined(t *testing.T) {
	errA := errors.New("error A")
	errB := errors.New("error B")
	a := &mockNotifier{err: errA}
	b := &mockNotifier{err: errB}
	tee, _ := NewTeeNotifier(a, b)

	err := tee.Send(context.Background(), teeMsg())
	if err == nil {
		t.Fatal("expected combined error")
	}
	if !errors.Is(err, errA) {
		t.Errorf("expected joined error to contain errA")
	}
	if !errors.Is(err, errB) {
		t.Errorf("expected joined error to contain errB")
	}
}
