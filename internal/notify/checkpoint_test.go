package notify

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewCheckpointNotifier_NilInner(t *testing.T) {
	_, err := NewCheckpointNotifier(nil)
	if err == nil {
		t.Fatal("expected error for nil inner, got nil")
	}
}

func TestNewCheckpointNotifier_Valid(t *testing.T) {
	cp, err := NewCheckpointNotifier(NewNoopNotifier())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cp == nil {
		t.Fatal("expected non-nil CheckpointNotifier")
	}
}

func TestCheckpointNotifier_LastSeen_Missing(t *testing.T) {
	cp, _ := NewCheckpointNotifier(NewNoopNotifier())
	_, ok := cp.LastSeen("secret/missing")
	if ok {
		t.Fatal("expected false for unseen path")
	}
}

func TestCheckpointNotifier_Send_RecordsSuccess(t *testing.T) {
	cp, _ := NewCheckpointNotifier(NewNoopNotifier())
	msg := Message{Path: "secret/db", Status: StatusExpiringSoon}

	if err := cp.Send(context.Background(), msg); err != nil {
		t.Fatalf("unexpected send error: %v", err)
	}

	rec, ok := cp.LastSeen("secret/db")
	if !ok {
		t.Fatal("expected record to exist after send")
	}
	if !rec.Succeeded {
		t.Error("expected Succeeded=true")
	}
	if rec.Status != StatusExpiringSoon {
		t.Errorf("expected status %v, got %v", StatusExpiringSoon, rec.Status)
	}
	if time.Since(rec.SentAt) > 2*time.Second {
		t.Error("SentAt is too old")
	}
}

func TestCheckpointNotifier_Send_RecordsFailure(t *testing.T) {
	sentinel := errors.New("inner failure")
	failing := &mockNotifier{err: sentinel}
	cp, _ := NewCheckpointNotifier(failing)

	msg := Message{Path: "secret/api", Status: StatusExpired}
	_ = cp.Send(context.Background(), msg)

	rec, ok := cp.LastSeen("secret/api")
	if !ok {
		t.Fatal("expected record even on failure")
	}
	if rec.Succeeded {
		t.Error("expected Succeeded=false on inner error")
	}
}

func TestCheckpointNotifier_All_ReturnsAllPaths(t *testing.T) {
	cp, _ := NewCheckpointNotifier(NewNoopNotifier())
	paths := []string{"secret/a", "secret/b", "secret/c"}
	for _, p := range paths {
		_ = cp.Send(context.Background(), Message{Path: p, Status: StatusOK})
	}

	all := cp.All()
	if len(all) != len(paths) {
		t.Errorf("expected %d records, got %d", len(paths), len(all))
	}
}
