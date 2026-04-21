package notify_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/notify"
)

func transformMsg() notify.Message {
	return notify.Message{
		Path:      "secret/data/db",
		Status:    notify.StatusExpiringSoon,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		Body:      "expiring soon",
	}
}

func TestNewTransformNotifier_NilInner(t *testing.T) {
	_, err := notify.NewTransformNotifier(nil, func(m notify.Message) notify.Message { return m })
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewTransformNotifier_NilFn(t *testing.T) {
	noop := notify.NewNoopNotifier()
	_, err := notify.NewTransformNotifier(noop, nil)
	if err == nil {
		t.Fatal("expected error for nil fn")
	}
}

func TestNewTransformNotifier_Valid(t *testing.T) {
	noop := notify.NewNoopNotifier()
	_, err := notify.NewTransformNotifier(noop, func(m notify.Message) notify.Message { return m })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTransformNotifier_AppliesTransform(t *testing.T) {
	var got notify.Message
	cap := &capturingNotifier{fn: func(m notify.Message) error { got = m; return nil }}

	tn, _ := notify.NewTransformNotifier(cap, func(m notify.Message) notify.Message {
		m.Body = strings.ToUpper(m.Body)
		return m
	})

	msg := transformMsg()
	if err := tn.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Body != "EXPIRING SOON" {
		t.Errorf("expected transformed body, got %q", got.Body)
	}
}

func TestTransformNotifier_ForwardsInnerError(t *testing.T) {
	sentinel := errors.New("inner error")
	cap := &capturingNotifier{fn: func(m notify.Message) error { return sentinel }}

	tn, _ := notify.NewTransformNotifier(cap, func(m notify.Message) notify.Message { return m })

	if err := tn.Send(transformMsg()); !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

// capturingNotifier records the last message sent.
type capturingNotifier struct {
	fn func(notify.Message) error
}

func (c *capturingNotifier) Send(m notify.Message) error { return c.fn(m) }
