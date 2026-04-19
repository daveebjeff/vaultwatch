package notify

import (
	"fmt"
	"sync"
	"time"
)

// AcknowledgeNotifier wraps a Notifier and suppresses re-delivery of alerts
// that have been explicitly acknowledged, until the acknowledgement expires.
type AcknowledgeNotifier struct {
	inner   Notifier
	mu      sync.Mutex
	acked   map[string]time.Time
	ttl     time.Duration
}

// NewAcknowledgeNotifier returns an AcknowledgeNotifier that suppresses
// forwarding for acknowledged secret paths for the given ttl duration.
func NewAcknowledgeNotifier(inner Notifier, ttl time.Duration) (*AcknowledgeNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("acknowledge: inner notifier must not be nil")
	}
	if ttl <= 0 {
		return nil, fmt.Errorf("acknowledge: ttl must be positive")
	}
	return &AcknowledgeNotifier{
		inner: inner,
		acked: make(map[string]time.Time),
		ttl:   ttl,
	}, nil
}

// Acknowledge marks the given secret path as acknowledged until ttl elapses.
func (a *AcknowledgeNotifier) Acknowledge(path string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.acked[path] = time.Now().Add(a.ttl)
}

// IsAcknowledged reports whether the path is currently acknowledged.
func (a *AcknowledgeNotifier) IsAcknowledged(path string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	expiry, ok := a.acked[path]
	if !ok {
		return false
	}
	if time.Now().After(expiry) {
		delete(a.acked, path)
		return false
	}
	return true
}

// Send forwards the message only if the secret path has not been acknowledged.
func (a *AcknowledgeNotifier) Send(msg Message) error {
	if a.IsAcknowledged(msg.Path) {
		return nil
	}
	return a.inner.Send(msg)
}
