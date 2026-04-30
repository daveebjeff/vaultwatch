package notify

import (
	"context"
	"sync"
	"time"
)

// ObserveNotifier wraps a Notifier and records latency and outcome for every
// Send call. Callers can retrieve a snapshot of collected metrics via Stats.
type ObserveNotifier struct {
	inner    Notifier
	mu       sync.Mutex
	total    int64
	errors   int64
	totalNs  int64 // cumulative nanoseconds
}

// ObserveStats is a point-in-time snapshot of collected observations.
type ObserveStats struct {
	Total      int64
	Errors     int64
	AvgLatency time.Duration
}

// NewObserveNotifier returns an ObserveNotifier that wraps inner.
// It returns an error when inner is nil.
func NewObserveNotifier(inner Notifier) (*ObserveNotifier, error) {
	if inner == nil {
		return nil, ErrNilInner
	}
	return &ObserveNotifier{inner: inner}, nil
}

// Send forwards the message to the inner notifier and records latency and
// whether the call succeeded.
func (o *ObserveNotifier) Send(ctx context.Context, msg Message) error {
	start := time.Now()
	err := o.inner.Send(ctx, msg)
	elapsed := time.Since(start)

	o.mu.Lock()
	o.total++
	o.totalNs += elapsed.Nanoseconds()
	if err != nil {
		o.errors++
	}
	o.mu.Unlock()

	return err
}

// Stats returns a snapshot of the collected metrics.
func (o *ObserveNotifier) Stats() ObserveStats {
	o.mu.Lock()
	defer o.mu.Unlock()

	var avg time.Duration
	if o.total > 0 {
		avg = time.Duration(o.totalNs / o.total)
	}
	return ObserveStats{
		Total:      o.total,
		Errors:     o.errors,
		AvgLatency: avg,
	}
}

// Reset clears all recorded observations.
func (o *ObserveNotifier) Reset() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.total = 0
	o.errors = 0
	o.totalNs = 0
}
