package notify

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// QuotaNotifier enforces a hard cap on the total number of notifications
// sent within a rolling time window, across all secret paths. Once the
// quota is exhausted the notifier returns ErrQuotaExceeded until the
// window resets.
type QuotaNotifier struct {
	inner    Notifier
	max      int
	window   time.Duration

	mu       sync.Mutex
	count    int
	windowAt time.Time
}

// NewQuotaNotifier wraps inner and allows at most max Send calls per window.
func NewQuotaNotifier(inner Notifier, max int, window time.Duration) (*QuotaNotifier, error) {
	if inner == nil {
		return nil, ErrNilInner
	}
	if max <= 0 {
		return nil, fmt.Errorf("quota: max must be > 0, got %d", max)
	}
	if window <= 0 {
		return nil, fmt.Errorf("quota: window must be > 0, got %s", window)
	}
	return &QuotaNotifier{
		inner:    inner,
		max:      max,
		window:   window,
		windowAt: time.Now(),
	}, nil
}

// Send forwards the message to the inner notifier if the quota has not
// been reached for the current window. Returns ErrQuotaExceeded otherwise.
func (q *QuotaNotifier) Send(ctx context.Context, msg Message) error {
	q.mu.Lock()
	now := time.Now()
	if now.After(q.windowAt.Add(q.window)) {
		q.count = 0
		q.windowAt = now
	}
	if q.count >= q.max {
		q.mu.Unlock()
		return ErrQuotaExceeded
	}
	q.count++
	q.mu.Unlock()
	return q.inner.Send(ctx, msg)
}

// Remaining returns the number of Send calls still permitted in the current
// window, along with the time at which the window resets.
func (q *QuotaNotifier) Remaining() (int, time.Time) {
	q.mu.Lock()
	defer q.mu.Unlock()
	now := time.Now()
	if now.After(q.windowAt.Add(q.window)) {
		return q.max, now.Add(q.window)
	}
	remaining := q.max - q.count
	if remaining < 0 {
		remaining = 0
	}
	return remaining, q.windowAt.Add(q.window)
}
