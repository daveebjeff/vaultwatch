package notify

import (
	"testing"
	"time"
)

// Note: syslog is only available on Unix; these tests verify construction and
// basic Send behaviour without mocking the kernel socket.

func TestNewSyslogNotifier_DefaultTag(t *testing.T) {
	n, err := NewSyslogNotifier("")
	if err != nil {
		t.Skipf("syslog unavailable: %v", err)
	}
	defer n.Close()
	if n.tag != "vaultwatch" {
		t.Errorf("expected default tag 'vaultwatch', got %q", n.tag)
	}
}

func TestNewSyslogNotifier_CustomTag(t *testing.T) {
	n, err := NewSyslogNotifier("myapp")
	if err != nil {
		t.Skipf("syslog unavailable: %v", err)
	}
	defer n.Close()
	if n.tag != "myapp" {
		t.Errorf("expected tag 'myapp', got %q", n.tag)
	}
}

func TestSyslogNotifier_Send_Expired(t *testing.T) {
	n, err := NewSyslogNotifier("vaultwatch-test")
	if err != nil {
		t.Skipf("syslog unavailable: %v", err)
	}
	defer n.Close()
	msg := Message{
		Status:     StatusExpired,
		SecretPath: "secret/db/prod",
		ExpiresAt:  time.Now().Add(-time.Hour),
	}
	if err := n.Send(msg); err != nil {
		t.Errorf("unexpected Send error: %v", err)
	}
}

func TestSyslogNotifier_Send_ExpiringSoon(t *testing.T) {
	n, err := NewSyslogNotifier("vaultwatch-test")
	if err != nil {
		t.Skipf("syslog unavailable: %v", err)
	}
	defer n.Close()
	msg := Message{
		Status:     StatusExpiringSoon,
		SecretPath: "secret/api/key",
		ExpiresAt:  time.Now().Add(30 * time.Minute),
	}
	if err := n.Send(msg); err != nil {
		t.Errorf("unexpected Send error: %v", err)
	}
}
