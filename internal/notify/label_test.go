package notify

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func labelMsg() Message {
	return Message{
		Path:      "secret/data/db",
		Status:    StatusExpiringSoon,
		Summary:   "expires soon",
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
}

func TestNewLabelNotifier_NilInner(t *testing.T) {
	_, err := NewLabelNotifier(nil, map[string]string{"env": "prod"})
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewLabelNotifier_NoLabels(t *testing.T) {
	_, err := NewLabelNotifier(NewNoopNotifier(), map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty labels")
	}
}

func TestNewLabelNotifier_Valid(t *testing.T) {
	n, err := NewLabelNotifier(NewNoopNotifier(), map[string]string{"env": "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestLabelNotifier_Send_PrefixesLabels(t *testing.T) {
	var got Message
	capture := &captureNotifier{fn: func(m Message) error { got = m; return nil }}
	n, _ := NewLabelNotifier(capture, map[string]string{"env": "staging"})

	n.Send(labelMsg())

	if !strings.Contains(got.Summary, "[env=staging]") {
		t.Errorf("expected label prefix in summary, got: %q", got.Summary)
	}
	if !strings.Contains(got.Summary, "expires soon") {
		t.Errorf("expected original summary preserved, got: %q", got.Summary)
	}
}

func TestLabelNotifier_Send_OriginalUnmodified(t *testing.T) {
	original := labelMsg()
	var got Message
	capture := &captureNotifier{fn: func(m Message) error { got = m; return nil }}
	n, _ := NewLabelNotifier(capture, map[string]string{"team": "ops"})

	n.Send(original)

	if got.Path != original.Path {
		t.Errorf("path changed: got %q want %q", got.Path, original.Path)
	}
	if got.Status != original.Status {
		t.Errorf("status changed")
	}
}

func TestLabelNotifier_Send_PropagatesError(t *testing.T) {
	fail := &captureNotifier{fn: func(m Message) error { return errors.New("send failed") }}
	n, _ := NewLabelNotifier(fail, map[string]string{"k": "v"})

	err := n.Send(labelMsg())
	if err == nil {
		t.Fatal("expected error to propagate")
	}
}

// captureNotifier is a test helper that calls fn on Send.
type captureNotifier struct {
	fn func(Message) error
}

func (c *captureNotifier) Send(m Message) error { return c.fn(m) }
