package notify

import (
	"fmt"
	"sync"
	"time"
)

// EscalationNotifier forwards to a primary notifier and, if it fails or the
// alert remains unacknowledged after a timeout, escalates to a secondary.
type EscalationNotifier struct {
	mu        sync.Mutex
	primary   Notifier
	secondary Notifier
	timeout   time.Duration
	pending   map[string]time.Time
}

// NewEscalationNotifier creates an EscalationNotifier.
// timeout is how long to wait before escalating an unacknowledged alert.
func NewEscalationNotifier(primary, secondary Notifier, timeout time.Duration) (*EscalationNotifier, error) {
	if primary == nil {
		return nil, fmt.Errorf("escalation: primary notifier must not be nil")
	}
	if secondary == nil {
		return nil, fmt.Errorf("escalation: secondary notifier must not be nil")
	}
	if timeout <= 0 {
		return nil, fmt.Errorf("escalation: timeout must be positive")
	}
	return &EscalationNotifier{
		primary:   primary,
		secondary: secondary,
		timeout:   timeout,
		pending:   make(map[string]time.Time),
	}, nil
}

// Send delivers via primary. If primary fails, secondary is tried immediately.
// On success, the alert is tracked; Escalate() checks for overdue ones.
func (e *EscalationNotifier) Send(msg Message) error {
	if err := e.primary.Send(msg); err != nil {
		return e.secondary.Send(msg)
	}
	e.mu.Lock()
	e.pending[msg.Path] = time.Now()
	e.mu.Unlock()
	return nil
}

// Acknowledge marks a secret path as handled, stopping escalation.
func (e *EscalationNotifier) Acknowledge(path string) {
	e.mu.Lock()
	delete(e.pending, path)
	e.mu.Unlock()
}

// Escalate checks all pending alerts and forwards overdue ones to secondary.
func (e *EscalationNotifier) Escalate(now time.Time) error {
	e.mu.Lock()
	overdue := make(map[string]time.Time)
	for path, sent := range e.pending {
		if now.Sub(sent) >= e.timeout {
			overdue[path] = sent
		}
	}
	e.mu.Unlock()

	var lastErr error
	for path := range overdue {
		msg := Message{
			Path:   path,
			Status: StatusExpired,
			Detail: "escalated: no acknowledgement within timeout",
		}
		if err := e.secondary.Send(msg); err != nil {
			lastErr = err
			continue
		}
		e.mu.Lock()
		delete(e.pending, path)
		e.mu.Unlock()
	}
	return lastErr
}
