// Package notify provides notification backends for vaultwatch alerts.
package notify

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo    Level = "INFO"
	LevelWarning Level = "WARNING"
	LevelCritical Level = "CRITICAL"
)

// Message holds the data for a single notification event.
type Message struct {
	Level     Level
	Secret    string
	ExpiresAt time.Time
	Details   string
}

// Notifier is the interface implemented by all notification backends.
type Notifier interface {
	Send(msg Message) error
}

// LogNotifier writes notifications as structured lines to a writer.
type LogNotifier struct {
	Out io.Writer
}

// NewLogNotifier returns a LogNotifier writing to stdout.
func NewLogNotifier() *LogNotifier {
	return &LogNotifier{Out: os.Stdout}
}

// Send writes the message to the configured writer.
func (l *LogNotifier) Send(msg Message) error {
	_, err := fmt.Fprintf(
		l.Out,
		"[%s] level=%s secret=%q expires_at=%s details=%q\n",
		time.Now().UTC().Format(time.RFC3339),
		msg.Level,
		msg.Secret,
		msg.ExpiresAt.UTC().Format(time.RFC3339),
		msg.Details,
	)
	return err
}
