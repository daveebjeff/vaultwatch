package notify

import (
	"errors"
	"testing"
)

func TestNewSnapshotNotifier_NilInner(t *testing.T) {
	_, err := NewSnapshotNotifier(nil)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewSnapshotNotifier_Valid(t *testing.T) {
	n, err := NewSnapshotNotifier(NewNoopNotifier())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestSnapshotNotifier_Latest_Missing(t *testing.T) {
	n, _ := NewSnapshotNotifier(NewNoopNotifier())
	_, ok := n.Latest("secret/foo")
	if ok {
		t.Fatal("expected no snapshot for unseen path")
	}
}

func TestSnapshotNotifier_Send_RecordsSnapshot(t *testing.T) {
	n, _ := NewSnapshotNotifier(NewNoopNotifier())
	msg := Message{Path: "secret/foo", Status: StatusExpiringSoon}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	snap, ok := n.Latest("secret/foo")
	if !ok {
		t.Fatal("expected snapshot after send")
	}
	if snap.Message.Status != StatusExpiringSoon {
		t.Errorf("expected ExpiringSoon, got %v", snap.Message.Status)
	}
	if snap.ReceivedAt.IsZero() {
		t.Error("expected non-zero ReceivedAt")
	}
}

func TestSnapshotNotifier_Send_InnerError_StillRecords(t *testing.T) {
	fail := &failNotifier{err: errors.New("boom")}
	n, _ := NewSnapshotNotifier(fail)
	msg := Message{Path: "secret/bar", Status: StatusExpired}
	err := n.Send(msg)
	if err == nil {
		t.Fatal("expected error from inner")
	}
	_, ok := n.Latest("secret/bar")
	if !ok {
		t.Fatal("expected snapshot even on inner error")
	}
}

func TestSnapshotNotifier_All_ReturnsCopy(t *testing.T) {
	n, _ := NewSnapshotNotifier(NewNoopNotifier())
	_ = n.Send(Message{Path: "a"})
	_ = n.Send(Message{Path: "b"})
	all := n.All()
	if len(all) != 2 {
		t.Errorf("expected 2 snapshots, got %d", len(all))
	}
}
