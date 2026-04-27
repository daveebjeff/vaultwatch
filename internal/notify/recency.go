package notify

import (
	"fmt"
	"sync"
	"time"
)

// RecencyNotifier suppresses repeated notifications for the same secret path
// unless the message status has changed or a minimum recency window has elapsed
// since the last forwarded notification.
//
// This is useful when the monitor loop fires frequently but you only want
// downstream notifiers to receive an alert once per recency window per path,
// unless the status changes (e.g. from ExpiringSoon to Expired).
type RecencyNotifier struct {
	inner    Notifier
	window   time.Duration
	mu       sync.Mutex
	lastSent map[string]lastEntry
}

type lastEntry struct {
	sentAt time.Time
	status Status
}

// NewRecencyNotifier wraps inner and suppresses duplicate notifications for
// the same path within window unless the status has changed.
func NewRecencyNotifier(inner Notifier, window time.Duration) (*RecencyNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("recency: inner notifier must not be nil")
	}
	if window <= 0 {
		return nil, fmt.Errorf("recency: window must be positive, got %s", window)
	}
	return &RecencyNotifier{
		inner:    inner,
		window:   window,
		lastSent: make(map[string]lastEntry),
	}, nil
}

// Send forwards msg to the inner notifier only if the path has not been
// notified within the recency window, or if the status has changed.
func (r *RecencyNotifier) Send(msg Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	entry, seen := r.lastSent[msg.Path]

	if seen && entry.status == msg.Status && now.Sub(entry.sentAt) < r.window {
		return nil
	}

	if err := r.inner.Send(msg); err != nil {
		return err
	}

	r.lastSent[msg.Path] = lastEntry{sentAt: now, status: msg.Status}
	return nil
}

// Reset clears all recency state, allowing the next send for every path
// to be forwarded unconditionally.
func (r *RecencyNotifier) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lastSent = make(map[string]lastEntry)
}
