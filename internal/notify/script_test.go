package notify

import (
	"runtime"
	"testing"
	"time"
)

func TestNewScriptNotifier_EmptyPath(t *testing.T) {
	_, err := NewScriptNotifier("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNewScriptNotifier_Valid(t *testing.T) {
	n, err := NewScriptNotifier("/usr/local/bin/alert.sh")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestScriptNotifier_Send_Success(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping on windows")
	}
	n, _ := NewScriptNotifier("/bin/true")
	msg := Message{
		Status:     StatusExpired,
		SecretPath: "secret/db",
		ExpiresAt:  time.Now(),
	}
	if err := n.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScriptNotifier_Send_Failure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping on windows")
	}
	n, _ := NewScriptNotifier("/bin/false")
	err := n.Send(Message{Status: StatusExpiringSoon, SecretPath: "secret/api"})
	if err == nil {
		t.Fatal("expected error when script exits non-zero")
	}
}

func TestScriptNotifier_Send_NotFound(t *testing.T) {
	n, _ := NewScriptNotifier("/nonexistent/script.sh")
	err := n.Send(Message{Status: StatusExpired, SecretPath: "secret/x"})
	if err == nil {
		t.Fatal("expected error for missing script")
	}
}
