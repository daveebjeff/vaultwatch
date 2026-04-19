package notify

import (
	"errors"
	"testing"
	"time"
)

func TestNewCircuitNotifier_NilInner(t *testing.T) {
	_, err := NewCircuitNotifier(nil, 3, time.Second)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewCircuitNotifier_ZeroMaxFailures(t *testing.T) {
	_, err := NewCircuitNotifier(NewNoopNotifier(), 0, time.Second)
	if err == nil {
		t.Fatal("expected error for zero maxFailures")
	}
}

func TestNewCircuitNotifier_ZeroReset(t *testing.T) {
	_, err := NewCircuitNotifier(NewNoopNotifier(), 3, 0)
	if err == nil {
		t.Fatal("expected error for zero resetAfter")
	}
}

func TestCircuitNotifier_ClosedOnSuccess(t *testing.T) {
	c, _ := NewCircuitNotifier(NewNoopNotifier(), 2, time.Second)
	if err := c.Send(exampleMsg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.State() != CircuitClosed {
		t.Errorf("expected closed, got %v", c.State())
	}
}

func TestCircuitNotifier_OpensAfterMaxFailures(t *testing.T) {
	fail := &mockFailNotifier{err: errors.New("boom")}
	c, _ := NewCircuitNotifier(fail, 2, time.Second)

	c.Send(exampleMsg) // failure 1
	c.Send(exampleMsg) // failure 2 — should open

	if c.State() != CircuitOpen {
		t.Errorf("expected open, got %v", c.State())
	}
}

func TestCircuitNotifier_OpenRejectsMessages(t *testing.T) {
	fail := &mockFailNotifier{err: errors.New("boom")}
	c, _ := NewCircuitNotifier(fail, 1, time.Hour)

	c.Send(exampleMsg) // opens circuit
	err := c.Send(exampleMsg)
	if err == nil {
		t.Fatal("expected error when circuit is open")
	}
}

func TestCircuitNotifier_HalfOpenAfterReset(t *testing.T) {
	fail := &mockFailNotifier{err: errors.New("boom")}
	c, _ := NewCircuitNotifier(fail, 1, time.Millisecond)

	c.Send(exampleMsg) // opens circuit
	time.Sleep(5 * time.Millisecond)

	// Next send should attempt (half-open), fail, and re-open
	c.Send(exampleMsg)
	if c.State() != CircuitOpen {
		t.Errorf("expected re-opened circuit, got %v", c.State())
	}
}

func TestCircuitNotifier_RecoveryClosesCircuit(t *testing.T) {
	fail := &mockFailNotifier{err: errors.New("boom")}
	c, _ := NewCircuitNotifier(fail, 1, time.Millisecond)

	c.Send(exampleMsg) // opens
	time.Sleep(5 * time.Millisecond)

	fail.err = nil // recover
	c.Send(exampleMsg)
	if c.State() != CircuitClosed {
		t.Errorf("expected closed after recovery, got %v", c.State())
	}
}

type mockFailNotifier struct{ err error }

func (m *mockFailNotifier) Send(_ Message) error { return m.err }
