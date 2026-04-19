package notify

import (
	"context"
	"math/rand"
	"time"
)

// JitterNotifier wraps a Notifier and introduces a random delay before
// forwarding each message. This helps avoid thundering-herd problems when
// many alerts fire simultaneously.
type JitterNotifier struct {
	inner   Notifier
	maxJitter time.Duration
	rng     *rand.Rand
}

// NewJitterNotifier returns a JitterNotifier that delays sends by a random
// duration in [0, maxJitter). inner must not be nil and maxJitter must be > 0.
func NewJitterNotifier(inner Notifier, maxJitter time.Duration) (*JitterNotifier, error) {
	if inner == nil {
		return nil, ErrNilInner
	}
	if maxJitter <= 0 {
		return nil, ErrInvalidConfig
	}
	return &JitterNotifier{
		inner:     inner,
		maxJitter: maxJitter,
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

// Send waits a random duration up to maxJitter, then forwards msg to the
// inner notifier. If ctx is cancelled during the wait, the send is skipped
// and the context error is returned.
func (j *JitterNotifier) Send(ctx context.Context, msg Message) error {
	delay := time.Duration(j.rng.Int63n(int64(j.maxJitter)))
	select {
	case <-time.After(delay):
		return j.inner.Send(ctx, msg)
	case <-ctx.Done():
		return ctx.Err()
	}
}
