package notify

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/yourusername/vaultwatch/internal/alert"
)

// EmailNotifier sends alert notifications via SMTP email.
type EmailNotifier struct {
	host     string
	port     int
	username string
	password string
	from     string
	to       []string
}

// NewEmailNotifier creates a new EmailNotifier.
// host, port, from, and at least one recipient are required.
func NewEmailNotifier(host string, port int, username, password, from string, to []string) (*EmailNotifier, error) {
	if host == "" {
		return nil, fmt.Errorf("email: SMTP host is required")
	}
	if from == "" {
		return nil, fmt.Errorf("email: sender address is required")
	}
	if len(to) == 0 {
		return nil, fmt.Errorf("email: at least one recipient is required")
	}
	return &EmailNotifier{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
		to:       to,
	}, nil
}

// Send delivers the alert message via SMTP.
func (e *EmailNotifier) Send(msg alert.Message) error {
	addr := fmt.Sprintf("%s:%d", e.host, e.port)
	subject := fmt.Sprintf("[VaultWatch] %s: %s", msg.Status, msg.Path)
	body := fmt.Sprintf(
		"Subject: %s\r\nFrom: %s\r\nTo: %s\r\n\r\nSecret: %s\nStatus: %s\nExpires: %s\n",
		subject, e.from, strings.Join(e.to, ", "),
		msg.Path, msg.Status, msg.Expiry.Format("2006-01-02 15:04:05 UTC"),
	)
	var auth smtp.Auth
	if e.username != "" {
		auth = smtp.PlainAuth("", e.username, e.password, e.host)
	}
	return smtp.SendMail(addr, auth, e.from, e.to, []byte(body))
}
