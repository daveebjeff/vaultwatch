package notify_test

import (
	"net"
	"net/smtp"
	"testing"
	"time"

	"github.com/yourusername/vaultwatch/internal/alert"
	"github.com/yourusername/vaultwatch/internal/notify"
)

func TestNewEmailNotifier_EmptyHost(t *testing.T) {
	_, err := notify.NewEmailNotifier("", 25, "", "", "from@example.com", []string{"to@example.com"})
	if err == nil {
		t.Fatal("expected error for empty host")
	}
}

func TestNewEmailNotifier_EmptyFrom(t *testing.T) {
	_, err := notify.NewEmailNotifier("smtp.example.com", 25, "", "", "", []string{"to@example.com"})
	if err == nil {
		t.Fatal("expected error for empty from")
	}
}

func TestNewEmailNotifier_NoRecipients(t *testing.T) {
	_, err := notify.NewEmailNotifier("smtp.example.com", 25, "", "", "from@example.com", nil)
	if err == nil {
		t.Fatal("expected error for no recipients")
	}
}

func TestNewEmailNotifier_Valid(t *testing.T) {
	n, err := notify.NewEmailNotifier("smtp.example.com", 587, "user", "pass", "from@example.com", []string{"to@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestEmailNotifier_Send_SMTPFailure(t *testing.T) {
	// Use a random available port with no real SMTP server to force failure.
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Skip("cannot bind local port")
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()

	_ = smtp.SendMail // ensure import used

	n, _ := notify.NewEmailNotifier("127.0.0.1", port, "", "", "from@example.com", []string{"to@example.com"})
	msg := alert.Message{
		Path:   "secret/db",
		Status: alert.StatusExpired,
		Expiry: time.Now().Add(-time.Hour),
	}
	err = n.Send(msg)
	if err == nil {
		t.Fatal("expected error when SMTP server unavailable")
	}
}
