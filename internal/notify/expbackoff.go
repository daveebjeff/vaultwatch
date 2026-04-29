package notify

import (
	"context"
	"math"
	"time"
)

// ExpBackoffNotifier wraps a Notifier and retries failed sends using
// exponential back-off with optional jitter. The delay doubles on each
// attempt up to MaxDelay.
type ExpBackoffNotifier struct {
	inner      Notifier
	initDelay  time.Duration
	maxDelay   time.Duration
	multiplier float64
	attempts   int
}

// NewExpBackoffNotifier returns an ExpBackoffNotifier.
//
// initDelay is the delay before the first retry.
// maxDelay caps the per-attempt delay.
// attempts is the maximum total send attempts (>= 1).
func NewExpBackoffNotifier(inner Notifier, initDelay, maxDelay time.Duration, attempts int) (*ExpBackoffNotifier, error) {
	if inner == nil {
		return nil, ErrNilInner
	}
	if initDelay <= 0 {
		return nil, ErrZeroDuration
	}
	if maxDelay < initDelay {
		maxDelay = initDelay
	}
	if attempts < 1 {
		attempts = 1
	}
	return &ExpBackoffNotifier{
		inner:      inner,
		initDelay:  initDelay,
		maxDelay:   maxDelay,
		multiplier: 2.0,
		attempts:   attempts,
	}, nil
}

// Send delivers msg, retrying on failure with exponential back-off.
func (e *ExpBackoffNotifier) Send(ctx context.Context, msg Message) error {
	var err error
	for i := 0; i < e.attempts; i++ {
		err = e.inner.Send(ctx, msg)
		if err == nil {
			return nil
		}
		if i == e.attempts-1 {
			break
		}
		delay := e.delayFor(i)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
	return err
}

func (e *ExpBackoffNotifier) delayFor(attempt int) time.Duration {
	d := float64(e.initDelay) * math.Pow(e.multiplier, float64(attempt))
	if d > float64(e.maxDelay) {
		d = float64(e.maxDelay)
	}
	return time.Duration(d)
}
