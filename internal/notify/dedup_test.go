package notify

import (
	"testing"
	"time"
)

func dedupMsg(path string, status Status) Message {
	return Message{Path: path, Status: status, ExpiresAt: time.Now().Add(time.Hour)}
}

func TestNewDedupNotifier_NilInner(t *testing.T) {
	_, err := NewDedupNotifier(nil)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestDedupNotifier_FirstSendAlwaysForwarded(t *testing.T) {
	cap := &captureNotifier{}
	d, _ := NewDedupNotifier(cap)
	d.Send(dedupMsg("secret/db", StatusExpiringSoon))
	if len(cap.msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(cap.msgs))
	}
}

func TestDedupNotifier_DuplicateSuppressed(t *testing.T) {
	cap := &captureNotifier{}
	d, _ := NewDedupNotifier(cap)
	d.Send(dedupMsg("secret/db", StatusExpiringSoon))
	d.Send(dedupMsg("secret/db", StatusExpiringSoon))
	if len(cap.msgs) != 1 {
		t.Fatalf("expected 1 message after duplicate, got %d", len(cap.msgs))
	}
}

func TestDedupNotifier_StatusChangeForwarded(t *testing.T) {
	cap := &captureNotifier{}
	d, _ := NewDedupNotifier(cap)
	d.Send(dedupMsg("secret/db", StatusExpiringSoon))
	d.Send(dedupMsg("secret/db", StatusExpired))
	if len(cap.msgs) != 2 {
		t.Fatalf("expected 2 messages on status change, got %d", len(cap.msgs))
	}
}

func TestDedupNotifier_IndependentPaths(t *testing.T) {
	cap := &captureNotifier{}
	d, _ := NewDedupNotifier(cap)
	d.Send(dedupMsg("secret/db", StatusExpiringSoon))
	d.Send(dedupMsg("secret/api", StatusExpiringSoon))
	if len(cap.msgs) != 2 {
		t.Fatalf("expected 2 messages for different paths, got %d", len(cap.msgs))
	}
}

func TestDedupNotifier_ResetAfterChange(t *testing.T) {
	cap := &captureNotifier{}
	d, _ := NewDedupNotifier(cap)
	d.Send(dedupMsg("secret/db", StatusExpiringSoon))
	d.Send(dedupMsg("secret/db", StatusExpired))
	d.Send(dedupMsg("secret/db", StatusExpired)) // duplicate of new state
	if len(cap.msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(cap.msgs))
	}
}
