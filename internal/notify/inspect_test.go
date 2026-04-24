package notify

import (
	"errors"
	"testing"
	"time"
)

func TestNewInspectNotifier_NilInner(t *testing.T) {
	_, err := NewInspectNotifier(nil)
	if err == nil {
		t.Fatal("expected error for nil inner, got nil")
	}
}

func TestNewInspectNotifier_Valid(t *testing.T) {
	noop := NewNoopNotifier()
	insp, err := NewInspectNotifier(noop)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if insp == nil {
		t.Fatal("expected non-nil InspectNotifier")
	}
}

func TestInspectNotifier_RecordsSuccess(t *testing.T) {
	noop := NewNoopNotifier()
	insp, _ := NewInspectNotifier(noop)

	msg := Message{Path: "secret/foo", Status: StatusExpiringSoon, ExpiresAt: time.Now().Add(time.Hour)}
	if err := insp.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if insp.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", insp.Len())
	}
	entries := insp.Entries()
	if entries[0].Message.Path != "secret/foo" {
		t.Errorf("expected path secret/foo, got %s", entries[0].Message.Path)
	}
	if entries[0].Err != nil {
		t.Errorf("expected nil error, got %v", entries[0].Err)
	}
}

func TestInspectNotifier_RecordsError(t *testing.T) {
	sentinel := errors.New("send failed")
	failing := &mockNotifier{err: sentinel}
	insp, _ := NewInspectNotifier(failing)

	msg := Message{Path: "secret/bar", Status: StatusExpired}
	_ = insp.Send(msg)

	entries := insp.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if !errors.Is(entries[0].Err, sentinel) {
		t.Errorf("expected sentinel error, got %v", entries[0].Err)
	}
}

func TestInspectNotifier_Reset(t *testing.T) {
	noop := NewNoopNotifier()
	insp, _ := NewInspectNotifier(noop)

	_ = insp.Send(Message{Path: "secret/a"})
	_ = insp.Send(Message{Path: "secret/b"})

	if insp.Len() != 2 {
		t.Fatalf("expected 2 entries before reset, got %d", insp.Len())
	}

	insp.Reset()

	if insp.Len() != 0 {
		t.Fatalf("expected 0 entries after reset, got %d", insp.Len())
	}
}

func TestInspectNotifier_EntriesAreSnapshot(t *testing.T) {
	noop := NewNoopNotifier()
	insp, _ := NewInspectNotifier(noop)

	_ = insp.Send(Message{Path: "secret/snap"})
	snap := insp.Entries()

	// Send another message after taking the snapshot.
	_ = insp.Send(Message{Path: "secret/snap2"})

	if len(snap) != 1 {
		t.Errorf("snapshot should be immutable; expected 1 entry, got %d", len(snap))
	}
}
