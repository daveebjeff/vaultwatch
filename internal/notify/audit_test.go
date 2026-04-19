package notify

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewAuditNotifier_NilInner(t *testing.T) {
	_, err := NewAuditNotifier(nil, "test", nil)
	if err == nil {
		t.Fatal("expected error for nil inner notifier")
	}
}

func TestNewAuditNotifier_Valid(t *testing.T) {
	n, err := NewAuditNotifier(NewNoopNotifier(), "noop", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestAuditNotifier_Send_WritesEntry(t *testing.T) {
	var buf bytes.Buffer
	n, _ := NewAuditNotifier(NewNoopNotifier(), "noop", &buf)

	msg := Message{
		Path:      "secret/db",
		Status:    StatusExpired,
		ExpiresAt: time.Now(),
	}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry AuditEntry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse audit entry: %v", err)
	}
	if entry.Path != "secret/db" {
		t.Errorf("expected path secret/db, got %s", entry.Path)
	}
	if entry.Status != StatusExpired {
		t.Errorf("expected status expired, got %v", entry.Status)
	}
	if entry.Notifier != "noop" {
		t.Errorf("expected notifier noop, got %s", entry.Notifier)
	}
	if entry.Error != "" {
		t.Errorf("expected no error field, got %s", entry.Error)
	}
}

func TestAuditNotifier_Send_RecordsError(t *testing.T) {
	var buf bytes.Buffer
	failing := &mockFailNotifier{err: errors.New("send failed")}
	n, _ := NewAuditNotifier(failing, "failing", &buf)

	msg := Message{Path: "secret/api", Status: StatusExpiringSoon}
	err := n.Send(msg)
	if err == nil {
		t.Fatal("expected error from inner notifier")
	}

	if !strings.Contains(buf.String(), "send failed") {
		t.Errorf("expected audit log to contain error text, got: %s", buf.String())
	}
}

func TestAuditNotifier_DefaultName(t *testing.T) {
	var buf bytes.Buffer
	n, _ := NewAuditNotifier(NewNoopNotifier(), "", &buf)
	_ = n.Send(Message{Path: "x", Status: StatusOK})

	var entry AuditEntry
	_ = json.Unmarshal(buf.Bytes(), &entry)
	if entry.Notifier != "unknown" {
		t.Errorf("expected notifier 'unknown', got %s", entry.Notifier)
	}
}

type mockFailNotifier struct{ err error }

func (m *mockFailNotifier) Send(_ Message) error { return m.err }
