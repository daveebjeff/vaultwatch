package notify

import (
	"context"
	"sync"
	"time"
)

// DeadLetterNotifier captures messages that failed to send for later inspection.
type DeadLetterNotifier struct {
	inner    Notifier
	mu       sync.Mutex
	failed   []DeadLetter
	maxSize  int
}

// DeadLetter holds a failed message and the error that caused the failure.
type DeadLetter struct {
	Message   Message
	Err       error
	FailedAt  time.Time
}

// NewDeadLetterNotifier wraps inner and retains up to maxSize failed messages.
func NewDeadLetterNotifier(inner Notifier, maxSize int) (*DeadLetterNotifier, error) {
	if inner == nil {
		return nil, ErrNilInner
	}
	if maxSize <= 0 {
		maxSize = 100
	}
	return &DeadLetterNotifier{inner: inner, maxSize: maxSize}, nil
}

// Send forwards the message and captures it on failure.
func (d *DeadLetterNotifier) Send(ctx context.Context, msg Message) error {
	err := d.inner.Send(ctx, msg)
	if err != nil {
		d.mu.Lock()
		defer d.mu.Unlock()
		if len(d.failed) < d.maxSize {
			d.failed = append(d.failed, DeadLetter{
				Message:  msg,
				Err:      err,
				FailedAt: time.Now(),
			})
		}
	}
	return err
}

// Failed returns a copy of all captured dead-letter entries.
func (d *DeadLetterNotifier) Failed() []DeadLetter {
	d.mu.Lock()
	defer d.mu.Unlock()
	out := make([]DeadLetter, len(d.failed))
	copy(out, d.failed)
	return out
}

// Drain returns and clears all captured dead-letter entries.
func (d *DeadLetterNotifier) Drain() []DeadLetter {
	d.mu.Lock()
	defer d.mu.Unlock()
	out := d.failed
	d.failed = nil
	return out
}
