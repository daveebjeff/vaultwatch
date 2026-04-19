package notify

import (
	"errors"
	"testing"
	"time"
)

func TestNewEscalationNotifier_NilPrimary(t *testing.T) {
	_, err := NewEscalationNotifier(nil, NewNoopNotifier(), time.Minute)
	if err == nil {
		t.Fatal("expected error for nil primary")
	}
}

func TestNewEscalationNotifier_NilSecondary(t *testing.T) {
	_, err := NewEscalationNotifier(NewNoopNotifier(), nil, time.Minute)
	if err == nil {
		t.Fatal("expected error for nil secondary")
	}
}

func TestNewEscalationNotifier_ZeroTimeout(t *testing.T) {
	_, err := NewEscalationNotifier(NewNoopNotifier(), NewNoopNotifier(), 0)
	if err == nil {
		t.Fatal("expected error for zero timeout")
	}
}

func TestEscalationNotifier_PrimarySuccess(t *testing.T) {
	primary := NewNoopNotifier()
	secondary := &recordingNotifier{}
	en, _ := NewEscalationNotifier(primary, secondary, time.Minute)

	msg := Message{Path: "secret/foo", Status: StatusExpired}
	if err := en.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secondary.count != 0 {
		t.Error("secondary should not be called when primary succeeds")
	}
}

func TestEscalationNotifier_PrimaryFailureFallsBack(t *testing.T) {
	primary := &failingNotifier{err: errors.New("primary down")}
	secondary := &recordingNotifier{}
	en, _ := NewEscalationNotifier(primary, secondary, time.Minute)

	msg := Message{Path: "secret/foo", Status: StatusExpired}
	if err := en.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secondary.count != 1 {
		t.Errorf("expected secondary called once, got %d", secondary.count)
	}
}

func TestEscalationNotifier_EscalatesAfterTimeout(t *testing.T) {
	primary := NewNoopNotifier()
	secondary := &recordingNotifier{}
	en, _ := NewEscalationNotifier(primary, secondary, 5*time.Minute)

	msg := Message{Path: "secret/db", Status: StatusExpiringSoon}
	en.Send(msg) //nolint

	// Not yet overdue
	en.Escalate(time.Now())
	if secondary.count != 0 {
		t.Error("should not escalate before timeout")
	}

	// Simulate timeout passed
	en.Escalate(time.Now().Add(10 * time.Minute))
	if secondary.count != 1 {
		t.Errorf("expected escalation, got %d secondary calls", secondary.count)
	}
}

func TestEscalationNotifier_AcknowledgePreventsEscalation(t *testing.T) {
	primary := NewNoopNotifier()
	secondary := &recordingNotifier{}
	en, _ := NewEscalationNotifier(primary, secondary, time.Minute)

	msg := Message{Path: "secret/ack", Status: StatusExpiringSoon}
	en.Send(msg) //nolint
	en.Acknowledge("secret/ack")
	en.Escalate(time.Now().Add(10 * time.Minute))

	if secondary.count != 0 {
		t.Error("acknowledged alert should not escalate")
	}
}

// helpers

type recordingNotifier struct{ count int }

func (r *recordingNotifier) Send(_ Message) error { r.count++; return nil }

type failingNotifier struct{ err error }

func (f *failingNotifier) Send(_ Message) error { return f.err }
