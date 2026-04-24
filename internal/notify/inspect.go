package notify

import (
	"fmt"
	"sync"
	"time"
)

// InspectEntry records a single notification event captured by an InspectNotifier.
type InspectEntry struct {
	Message   Message
	SentAt    time.Time
	Err       error
}

// InspectNotifier wraps an inner Notifier and records every Send call,
// including the message and any error returned. It is primarily intended
// for use in tests and debugging pipelines.
type InspectNotifier struct {
	inner   Notifier
	mu      sync.Mutex
	entries []InspectEntry
}

// NewInspectNotifier returns an InspectNotifier wrapping inner.
// inner must not be nil.
func NewInspectNotifier(inner Notifier) (*InspectNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("inspect: inner notifier must not be nil")
	}
	return &InspectNotifier{inner: inner}, nil
}

// Send forwards the message to the inner notifier and records the result.
func (n *InspectNotifier) Send(msg Message) error {
	err := n.inner.Send(msg)
	n.mu.Lock()
	n.entries = append(n.entries, InspectEntry{
		Message: msg,
		SentAt:  time.Now(),
		Err:     err,
	})
	n.mu.Unlock()
	return err
}

// Entries returns a snapshot of all recorded entries in order of receipt.
func (n *InspectNotifier) Entries() []InspectEntry {
	n.mu.Lock()
	defer n.mu.Unlock()
	out := make([]InspectEntry, len(n.entries))
	copy(out, n.entries)
	return out
}

// Len returns the number of Send calls recorded so far.
func (n *InspectNotifier) Len() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return len(n.entries)
}

// Reset clears all recorded entries.
func (n *InspectNotifier) Reset() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.entries = nil
}
