package notify

import (
	"errors"
	"testing"
	"time"
)

func TestNewEnrichNotifier_NilInner(t *testing.T) {
	_, err := NewEnrichNotifier(nil)
	if err == nil {
		t.Fatal("expected error for nil inner, got nil")
	}
}

func TestNewEnrichNotifier_Valid(t *testing.T) {
	n, err := NewEnrichNotifier(NewNoopNotifier())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestEnrichNotifier_SeverityCritical(t *testing.T) {
	cap := &capturingNotifier{}
	en, _ := NewEnrichNotifier(cap)

	msg := Message{Path: "secret/db", Status: StatusExpired, Expiry: time.Now().Add(-time.Hour)}
	if err := en.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := cap.last.Labels["severity"]; got != "critical" {
		t.Errorf("severity = %q, want %q", got, "critical")
	}
	if got := cap.last.Labels["time_to_expiry"]; got != "expired" {
		t.Errorf("time_to_expiry = %q, want \"expired\"", got)
	}
}

func TestEnrichNotifier_SeverityWarning(t *testing.T) {
	cap := &capturingNotifier{}
	en, _ := NewEnrichNotifier(cap)

	future := time.Now().Add(2 * time.Hour)
	msg := Message{Path: "secret/db", Status: StatusExpiringSoon, Expiry: future}
	en.now = func() time.Time { return future.Add(-2 * time.Hour) }

	if err := en.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := cap.last.Labels["severity"]; got != "warning" {
		t.Errorf("severity = %q, want %q", got, "warning")
	}
	if cap.last.Labels["time_to_expiry"] == "" {
		t.Error("expected non-empty time_to_expiry label")
	}
}

func TestEnrichNotifier_SeverityInfo(t *testing.T) {
	cap := &capturingNotifier{}
	en, _ := NewEnrichNotifier(cap)

	msg := Message{Path: "secret/ok", Status: StatusOK}
	if err := en.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := cap.last.Labels["severity"]; got != "info" {
		t.Errorf("severity = %q, want %q", got, "info")
	}
}

func TestEnrichNotifier_ZeroExpiry_NoTTLLabel(t *testing.T) {
	cap := &capturingNotifier{}
	en, _ := NewEnrichNotifier(cap)

	msg := Message{Path: "secret/ok", Status: StatusOK} // zero Expiry
	en.Send(msg)

	if _, ok := cap.last.Labels["time_to_expiry"]; ok {
		t.Error("time_to_expiry label should be absent when Expiry is zero")
	}
}

func TestEnrichNotifier_PreservesExistingLabels(t *testing.T) {
	cap := &capturingNotifier{}
	en, _ := NewEnrichNotifier(cap)

	msg := Message{
		Path:   "secret/db",
		Status: StatusOK,
		Labels: map[string]string{"env": "prod"},
	}
	en.Send(msg)

	if got := cap.last.Labels["env"]; got != "prod" {
		t.Errorf("existing label env = %q, want %q", got, "prod")
	}
}

func TestEnrichNotifier_DoesNotMutateOriginal(t *testing.T) {
	cap := &capturingNotifier{}
	en, _ := NewEnrichNotifier(cap)

	msg := Message{Path: "secret/db", Status: StatusExpired}
	en.Send(msg)

	if msg.Labels != nil {
		t.Error("original message Labels should remain nil")
	}
}

func TestEnrichNotifier_PropagatesInnerError(t *testing.T) {
	sentinel := errors.New("inner failure")
	inner := &errorNotifier{err: sentinel}
	en, _ := NewEnrichNotifier(inner)

	if err := en.Send(Message{Path: "x", Status: StatusOK}); !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

// capturingNotifier records the last message it received.
type capturingNotifier struct{ last Message }

func (c *capturingNotifier) Send(m Message) error { c.last = m; return nil }

// errorNotifier always returns the configured error.
type errorNotifier struct{ err error }

func (e *errorNotifier) Send(_ Message) error { return e.err }
