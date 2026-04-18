package notify

import (
	"fmt"
	"log/syslog"
)

// SyslogNotifier sends alert messages to the local syslog daemon.
type SyslogNotifier struct {
	writer *syslog.Writer
	tag    string
}

// NewSyslogNotifier creates a SyslogNotifier writing under the given tag.
// Priority is LOG_WARNING | LOG_DAEMON.
func NewSyslogNotifier(tag string) (*SyslogNotifier, error) {
	if tag == "" {
		tag = "vaultwatch"
	}
	w, err := syslog.New(syslog.LOG_WARNING|syslog.LOG_DAEMON, tag)
	if err != nil {
		return nil, fmt.Errorf("syslog: open: %w", err)
	}
	return &SyslogNotifier{writer: w, tag: tag}, nil
}

// Send writes the alert message to syslog, choosing severity by status.
func (s *SyslogNotifier) Send(msg Message) error {
	text := fmt.Sprintf("[%s] secret=%s expires=%s",
		msg.Status, msg.SecretPath, msg.ExpiresAt.Format("2006-01-02T15:04:05Z"))

	switch msg.Status {
	case StatusExpired:
		return s.writer.Err(text)
	case StatusExpiringSoon:
		return s.writer.Warning(text)
	default:
		return s.writer.Info(text)
	}
}

// Close releases the underlying syslog connection.
func (s *SyslogNotifier) Close() error {
	return s.writer.Close()
}
