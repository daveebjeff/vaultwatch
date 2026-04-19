package notify

import (
	"context"
	"fmt"
	"time"
)

// TimeoutNotifier wraps a Notifier and enforces a maximum send duration.
// If the inner Send does not complete within the deadline, it is cancelled
// and an error is returned.
type TimeoutNotifier struct {
	inner   Notifier
	timeout time.Duration
}

// NewTimeoutNotifier returns a TimeoutNotifier that cancels sends exceeding d.
func NewTimeoutNotifier(inner Notifier, d time.Duration) (*TimeoutNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("timeout: inner notifier must not be nil")
	}
	if d <= 0 {
		return nil, fmt.Errorf("timeout: duration must be positive")
	}
	return &TimeoutNotifier{inner: inner, timeout: d}, nil
}

// Send forwards msg to the inner notifier, returning an error if the timeout
// elapses before the send completes.
func (t *TimeoutNotifier) Send(ctx context.Context, msg Message) error {
	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

	type result struct{ err error }
	ch := make(chan result, 1)
	go func() {
		ch <- result{err: t.inner.Send(ctx, msg)}
	}()

	select {
	case r := <-ch:
		return r.err
	case <-ctx.Done():
		return fmt.Errorf("timeout: send exceeded %s: %w", t.timeout, ctx.Err())
	}
}
