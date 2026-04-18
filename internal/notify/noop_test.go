package notify_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultwatch/internal/notify"
)

func TestNewNoopNotifier_NotNil(t *testing.T) {
	n := notify.NewNoopNotifier()
	if n == nil {
		t.Fatal("expected non-nil NoopNotifier")
	}
}

func TestNoopNotifier_Send_ReturnsNil(t *testing.T) {
	n := notify.NewNoopNotifier()
	msg := notify.Message{
		Path:      "secret/my-app/db",
		Status:    notify.StatusExpired,
		ExpiresAt: time.Now().Add(-time.Minute),
	}
	if err := n.Send(msg); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestNoopNotifier_Send_MultipleCallsOK(t *testing.T) {
	n := notify.NewNoopNotifier()
	for i := 0; i < 10; i++ {
		msg := notify.Message{
			Path:   "secret/path",
			Status: notify.StatusExpiringSoon,
		}
		if err := n.Send(msg); err != nil {
			t.Fatalf("call %d: expected nil error, got %v", i, err)
		}
	}
}

func TestNoopNotifier_ImplementsNotifier(t *testing.T) {
	var _ interface{ Send(notify.Message) error } = notify.NewNoopNotifier()
}
