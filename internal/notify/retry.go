package notify

import (
	"fmt"
	"time"
)

// RetryNotifier wraps a Notifier and retries on failure up to MaxAttempts times.
type RetryNotifier struct {
	inner      Notifier
	maxAttempts int
	delay      time.Duration
}

// NewRetryNotifier creates a RetryNotifier wrapping the given Notifier.
// maxAttempts must be >= 1; delay is the wait between attempts.
func NewRetryNotifier(n Notifier, maxAttempts int, delay time.Duration) (*RetryNotifier, error) {
	if n == nil {
		return nil, fmt.Errorf("notify: inner notifier must not be nil")
	}
	if maxAttempts < 1 {
		return nil, fmt.Errorf("notify: maxAttempts must be at least 1")
	}
	return &RetryNotifier{
		inner:      n,
		maxAttempts: maxAttempts,
		delay:      delay,
	}, nil
}

// Send attempts to deliver the message, retrying up to maxAttempts times.
// Returns the last error if all attempts fail.
func (r *RetryNotifier) Send(msg Message) error {
	var err error
	for i := 0; i < r.maxAttempts; i++ {
		if err = r.inner.Send(msg); err == nil {
			return nil
		}
		if i < r.maxAttempts-1 && r.delay > 0 {
			time.Sleep(r.delay)
		}
	}
	return fmt.Errorf("notify: all %d attempts failed: %w", r.maxAttempts, err)
}
