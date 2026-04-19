package notify

import (
	"fmt"
	"sync"
	"time"
)

// ThrottleNotifier limits how many notifications can be sent in a given window.
// Once the limit is reached, additional messages are dropped until the window resets.
type ThrottleNotifier struct {
	inner    Notifier
	window   time.Duration
	maxCount int

	mu       sync.Mutex
	windowStart time.Time
	count    int
}

// NewThrottleNotifier wraps inner, allowing at most maxCount Send calls per window.
func NewThrottleNotifier(inner Notifier, maxCount int, window time.Duration) (*ThrottleNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("throttle: inner notifier must not be nil")
	}
	if maxCount <= 0 {
		return nil, fmt.Errorf("throttle: maxCount must be greater than zero")
	}
	if window <= 0 {
		return nil, fmt.Errorf("throttle: window must be greater than zero")
	}
	return &ThrottleNotifier{
		inner:       inner,
		window:      window,
		maxCount:    maxCount,
		windowStart: time.Now(),
	}, nil
}

// Send forwards the message to the inner notifier if the rate limit has not been exceeded.
func (t *ThrottleNotifier) Send(msg Message) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	if now.Sub(t.windowStart) >= t.window {
		t.windowStart = now
		t.count = 0
	}

	if t.count >= t.maxCount {
		return nil // silently drop
	}

	t.count++
	return t.inner.Send(msg)
}
