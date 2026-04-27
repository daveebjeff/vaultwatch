package notify

import (
	"context"
	"sync"
	"time"
)

// HedgeNotifier sends a notification to the primary notifier and, if it does
// not complete within the hedge delay, concurrently fires the secondary as a
// hedge. The first successful result wins; the other goroutine is abandoned.
type HedgeNotifier struct {
	primary   Notifier
	secondary Notifier
	delay     time.Duration
}

// NewHedgeNotifier returns a HedgeNotifier that hedges with secondary after
// delay. Both primary and secondary must be non-nil and delay must be > 0.
func NewHedgeNotifier(primary, secondary Notifier, delay time.Duration) (*HedgeNotifier, error) {
	if primary == nil {
		return nil, ErrNilInner
	}
	if secondary == nil {
		return nil, errNilSecondary
	}
	if delay <= 0 {
		return nil, errZeroDuration
	}
	return &HedgeNotifier{primary: primary, secondary: secondary, delay: delay}, nil
}

// Send dispatches msg to primary. If primary has not returned within the hedge
// delay, secondary is also invoked concurrently. The error from whichever
// notifier succeeds first (nil error) is returned. If both fail, the primary
// error is returned.
func (h *HedgeNotifier) Send(ctx context.Context, msg Message) error {
	type result struct {
		err error
	}

	primaryCh := make(chan result, 1)
	secondaryCh := make(chan result, 1)

	go func() {
		primaryCh <- result{err: h.primary.Send(ctx, msg)}
	}()

	timer := time.NewTimer(h.delay)
	defer timer.Stop()

	var once sync.Once
	launchSecondary := func() {
		once.Do(func() {
			go func() {
				secondaryCh <- result{err: h.secondary.Send(ctx, msg)}
			}()
		})
	}

	var primaryErr error
	primaryDone := false

	for {
		select {
		case <-timer.C:
			launchSecondary()
		case r := <-primaryCh:
			primaryErr = r.err
			primaryDone = true
			if r.err == nil {
				return nil
			}
			launchSecondary()
		case r := <-secondaryCh:
			if r.err == nil {
				return nil
			}
			if primaryDone {
				return primaryErr
			}
		case <-ctx.Done():
			return ctx.Err()
		}
		if primaryDone {
			// wait for secondary if it was launched
			select {
			case r := <-secondaryCh:
				if r.err == nil {
					return nil
				}
				return primaryErr
			case <-ctx.Done():
				return primaryErr
			default:
				// secondary not yet launched or not yet done; return primary result
				return primaryErr
			}
		}
	}
}
