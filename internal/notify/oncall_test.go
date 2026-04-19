package notify

import (
	"errors"
	"testing"
	"time"
)

func baseRotation(n Notifier, start, end time.Time) OnCallRotation {
	return OnCallRotation{Name: "test", Start: start, End: end, Notifier: n}
}

func TestNewOnCallNotifier_NoRotations(t *testing.T) {
	_, err := NewOnCallNotifier(nil)
	if err == nil {
		t.Fatal("expected error for empty rotations")
	}
}

func TestNewOnCallNotifier_NilNotifier(t *testing.T) {
	now := time.Now().UTC()
	_, err := NewOnCallNotifier([]OnCallRotation{
		{Name: "x", Start: now, End: now.Add(time.Hour), Notifier: nil},
	})
	if err == nil {
		t.Fatal("expected error for nil notifier")
	}
}

func TestNewOnCallNotifier_InvalidWindow(t *testing.T) {
	now := time.Now().UTC()
	_, err := NewOnCallNotifier([]OnCallRotation{
		baseRotation(NewNoopNotifier(), now.Add(time.Hour), now),
	})
	if err == nil {
		t.Fatal("expected error for end before start")
	}
}

func TestOnCallNotifier_ActiveRotationCalled(t *testing.T) {
	mock := &mockNotifier{}
	now := time.Now().UTC()
	on, err := NewOnCallNotifier([]OnCallRotation{
		baseRotation(mock, now.Add(-time.Hour), now.Add(time.Hour)),
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := on.Send(Message{Path: "secret/a"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.calls != 1 {
		t.Fatalf("expected 1 call, got %d", mock.calls)
	}
}

func TestOnCallNotifier_NoActiveRotation(t *testing.T) {
	now := time.Now().UTC()
	on, _ := NewOnCallNotifier([]OnCallRotation{
		baseRotation(NewNoopNotifier(), now.Add(-2*time.Hour), now.Add(-time.Hour)),
	})
	err := on.Send(Message{Path: "secret/b"})
	if !errors.Is(err, ErrNoOnCallRotation) {
		t.Fatalf("expected ErrNoOnCallRotation, got %v", err)
	}
}

func TestOnCallNotifier_AddRotation(t *testing.T) {
	mock := &mockNotifier{}
	now := time.Now().UTC()
	on, _ := NewOnCallNotifier([]OnCallRotation{
		baseRotation(NewNoopNotifier(), now.Add(-2*time.Hour), now.Add(-time.Hour)),
	})
	_ = on.AddRotation(baseRotation(mock, now.Add(-time.Minute), now.Add(time.Hour)))
	if err := on.Send(Message{Path: "secret/c"}); err != nil {
		t.Fatalf("unexpected error after AddRotation: %v", err)
	}
	if mock.calls != 1 {
		t.Fatalf("expected 1 call, got %d", mock.calls)
	}
}
