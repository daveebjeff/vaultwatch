package notify

import (
	"fmt"
	"sync"
	"time"
)

// RateLimitNotifier wraps a Notifier and suppresses duplicate alerts
// for the same secret path within a cooldown window.
type RateLimitNotifier struct {
	inner    Notifier
	cooldown time.Duration
	mu       sync.Mutex
	lastSent map[string]time.Time
}

// NewRateLimitNotifier creates a RateLimitNotifier with the given cooldown.
func NewRateLimitNotifier(n Notifier, cooldown time.Duration) (*RateLimitNotifier, error) {
	if n == nil {
		return nil, fmt.Errorf("notifier must not be nil")
	}
	if cooldown <= 0 {
		return nil, fmt.Errorf("cooldown must be positive")
	}
	return &RateLimitNotifier{
		inner:    n,
		cooldown: cooldown,
		lastSent: make(map[string]time.Time),
	}, nil
}

// Send forwards the message only if the cooldown has elapsed since the last
// notification for the same secret path.
func (r *RateLimitNotifier) Send(msg Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if last, ok := r.lastSent[msg.Path]; ok {
		if time.Since(last) < r.cooldown {
			return nil // suppressed
		}
	}
	r.lastSent[msg.Path] = time.Now()
	return r.inner.Send(msg)
}
