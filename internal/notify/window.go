package notify

import (
	"fmt"
	"sync"
	"time"
)

// WindowNotifier forwards messages only when the count of sends within a
// sliding time window stays below a configured ceiling. Unlike ThrottleNotifier
// (which resets on a fixed period) this uses a true sliding window.
type WindowNotifier struct {
	inner    Notifier
	max      int
	window   time.Duration
	mu       sync.Mutex
	timestamps []time.Time
}

// NewWindowNotifier creates a WindowNotifier that allows at most max sends
// within any rolling window of the given duration.
func NewWindowNotifier(inner Notifier, max int, window time.Duration) (*WindowNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("window: inner notifier must not be nil")
	}
	if max <= 0 {
		return nil, fmt.Errorf("window: max must be greater than zero")
	}
	if window <= 0 {
		return nil, fmt.Errorf("window: window duration must be greater than zero")
	}
	return &WindowNotifier{
		inner:  inner,
		max:    max,
		window: window,
	}, nil
}

// Send forwards the message to the inner notifier if the sliding-window limit
// has not been reached; otherwise it returns ErrSuppressed.
func (w *WindowNotifier) Send(msg Message) error {
	now := time.Now()
	cutoff := now.Add(-w.window)

	w.mu.Lock()
	// evict timestamps outside the window
	valid := w.timestamps[:0]
	for _, t := range w.timestamps {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	w.timestamps = valid

	if len(w.timestamps) >= w.max {
		w.mu.Unlock()
		return ErrSuppressed
	}
	w.timestamps = append(w.timestamps, now)
	w.mu.Unlock()

	return w.inner.Send(msg)
}
