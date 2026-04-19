package notify

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// AuditEntry records a single notification attempt.
type AuditEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Path      string    `json:"path"`
	Status    Status    `json:"status"`
	Notifier  string    `json:"notifier"`
	Error     string    `json:"error,omitempty"`
}

// AuditNotifier wraps another Notifier and writes a JSON audit log entry
// for every Send call.
type AuditNotifier struct {
	inner    Notifier
	name     string
	writer   io.Writer
}

// NewAuditNotifier returns an AuditNotifier that logs to w.
// name identifies the inner notifier in the audit log.
func NewAuditNotifier(inner Notifier, name string, w io.Writer) (*AuditNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("audit: inner notifier must not be nil")
	}
	if name == "" {
		name = "unknown"
	}
	if w == nil {
		w = os.Stderr
	}
	return &AuditNotifier{inner: inner, name: name, writer: w}, nil
}

// Send forwards the message to the inner notifier and writes an audit entry.
func (a *AuditNotifier) Send(msg Message) error {
	err := a.inner.Send(msg)

	entry := AuditEntry{
		Timestamp: time.Now().UTC(),
		Path:      msg.Path,
		Status:    msg.Status,
		Notifier:  a.name,
	}
	if err != nil {
		entry.Error = err.Error()
	}

	data, jsonErr := json.Marshal(entry)
	if jsonErr == nil {
		_, _ = fmt.Fprintf(a.writer, "%s\n", data)
	}

	return err
}
