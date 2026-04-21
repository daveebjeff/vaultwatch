package notify

import (
	"errors"
	"testing"
	"time"
)

func metaMsg() Message {
	return Message{
		Path:      "secret/db",
		Status:    StatusExpiringSoon,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
}

func TestNewMetadataNotifier_NilInner(t *testing.T) {
	_, err := NewMetadataNotifier(nil, map[string]string{"env": "prod"})
	if err == nil {
		t.Fatal("expected error for nil inner notifier")
	}
}

func TestNewMetadataNotifier_NoMetadata(t *testing.T) {
	_, err := NewMetadataNotifier(NewNoopNotifier(), map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty metadata")
	}
}

func TestNewMetadataNotifier_Valid(t *testing.T) {
	n, err := NewMetadataNotifier(NewNoopNotifier(), map[string]string{"env": "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestMetadataNotifier_StampsLabels(t *testing.T) {
	var got Message
	capture := &captureNotifier{fn: func(m Message) error { got = m; return nil }}

	n, _ := NewMetadataNotifier(capture, map[string]string{"region": "us-east-1", "env": "prod"})
	if err := n.Send(metaMsg()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Labels["region"] != "us-east-1" {
		t.Errorf("expected region label, got %v", got.Labels)
	}
	if got.Labels["env"] != "prod" {
		t.Errorf("expected env label, got %v", got.Labels)
	}
}

func TestMetadataNotifier_DoesNotOverwriteExisting(t *testing.T) {
	var got Message
	capture := &captureNotifier{fn: func(m Message) error { got = m; return nil }}

	n, _ := NewMetadataNotifier(capture, map[string]string{"env": "prod"})

	msg := metaMsg()
	msg.Labels = map[string]string{"env": "staging"}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Labels["env"] != "staging" {
		t.Errorf("expected existing label to be preserved, got %q", got.Labels["env"])
	}
}

func TestMetadataNotifier_InnerError(t *testing.T) {
	fail := &captureNotifier{fn: func(m Message) error { return errors.New("boom") }}
	n, _ := NewMetadataNotifier(fail, map[string]string{"k": "v"})
	if err := n.Send(metaMsg()); err == nil {
		t.Fatal("expected error from inner notifier")
	}
}

func TestMetadataNotifier_SetMetadata(t *testing.T) {
	var got Message
	capture := &captureNotifier{fn: func(m Message) error { got = m; return nil }}

	n, _ := NewMetadataNotifier(capture, map[string]string{"env": "prod"})
	if err := n.SetMetadata(map[string]string{"env": "dev", "team": "platform"}); err != nil {
		t.Fatalf("SetMetadata error: %v", err)
	}
	n.Send(metaMsg())
	if got.Labels["team"] != "platform" {
		t.Errorf("expected updated metadata, got %v", got.Labels)
	}
}

func TestMetadataNotifier_SetMetadata_Empty(t *testing.T) {
	n, _ := NewMetadataNotifier(NewNoopNotifier(), map[string]string{"k": "v"})
	if err := n.SetMetadata(map[string]string{}); err == nil {
		t.Fatal("expected error for empty replacement metadata")
	}
}

// captureNotifier is a test helper that calls fn for each Send.
type captureNotifier struct {
	fn func(Message) error
}

func (c *captureNotifier) Send(m Message) error { return c.fn(m) }
