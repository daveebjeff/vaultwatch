package notify

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewTraceIDNotifier_NilInner(t *testing.T) {
	_, err := NewTraceIDNotifier(nil, "")
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewTraceIDNotifier_DefaultHeader(t *testing.T) {
	n, err := NewTraceIDNotifier(NewNoopNotifier(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.header != "trace_id" {
		t.Errorf("expected header 'trace_id', got %q", n.header)
	}
}

func TestNewTraceIDNotifier_CustomHeader(t *testing.T) {
	n, err := NewTraceIDNotifier(NewNoopNotifier(), "x-request-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.header != "x-request-id" {
		t.Errorf("expected header 'x-request-id', got %q", n.header)
	}
}

func TestTraceIDNotifier_StampsLabel(t *testing.T) {
	cap := &capturingNotifier{}
	n, _ := NewTraceIDNotifier(cap, "")

	msg := Message{Path: "secret/foo", Status: StatusExpiringSoon, ExpiresAt: time.Now().Add(time.Hour)}
	if err := n.Send(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := cap.last.Labels["trace_id"]
	if got == "" {
		t.Error("expected trace_id label to be set")
	}
}

func TestTraceIDNotifier_ReusesContextTraceID(t *testing.T) {
	cap := &capturingNotifier{}
	n, _ := NewTraceIDNotifier(cap, "")

	ctx := ContextWithTraceID(context.Background(), "fixed-id-42")
	msg := Message{Path: "secret/bar", Status: StatusExpired, ExpiresAt: time.Now()}
	if err := n.Send(ctx, msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := cap.last.Labels["trace_id"]; got != "fixed-id-42" {
		t.Errorf("expected 'fixed-id-42', got %q", got)
	}
}

func TestTraceIDNotifier_UniquePerSend(t *testing.T) {
	cap := &capturingNotifier{}
	n, _ := NewTraceIDNotifier(cap, "")

	seen := make(map[string]bool)
	for i := 0; i < 10; i++ {
		msg := Message{Path: fmt.Sprintf("secret/%d", i), Status: StatusExpiringSoon, ExpiresAt: time.Now().Add(time.Hour)}
		_ = n.Send(context.Background(), msg)
		id := cap.last.Labels["trace_id"]
		if seen[id] {
			t.Errorf("duplicate trace ID generated: %s", id)
		}
		seen[id] = true
	}
}

func TestTraceIDNotifier_PreservesExistingLabels(t *testing.T) {
	cap := &capturingNotifier{}
	n, _ := NewTraceIDNotifier(cap, "")

	msg := Message{
		Path:      "secret/baz",
		Status:    StatusExpiringSoon,
		ExpiresAt: time.Now().Add(time.Hour),
		Labels:    map[string]string{"env": "prod"},
	}
	_ = n.Send(context.Background(), msg)

	if cap.last.Labels["env"] != "prod" {
		t.Error("existing label 'env' was overwritten")
	}
	if cap.last.Labels["trace_id"] == "" {
		t.Error("trace_id not set")
	}
}

// capturingNotifier records the last message it received.
type capturingNotifier struct {
	last Message
}

func (c *capturingNotifier) Send(_ context.Context, msg Message) error {
	c.last = msg
	return nil
}
