package notify

import (
	"context"
	"sync"
	"time"
)

// BatchNotifier collects messages over a window and sends them together
// as a slice to the inner notifier using a single summary message.
type BatchNotifier struct {
	inner    Notifier
	window   time.Duration
	maxSize  int
	mu       sync.Mutex
	batch    []Message
	timer    *time.Timer
}

// NewBatchNotifier creates a BatchNotifier that flushes after window duration
// or when maxSize messages have accumulated, whichever comes first.
func NewBatchNotifier(inner Notifier, window time.Duration, maxSize int) (*BatchNotifier, error) {
	if inner == nil {
		return nil, ErrNilInner
	}
	if window <= 0 {
		return nil, ErrZeroWindow
	}
	if maxSize <= 0 {
		maxSize = 50
	}
	return &BatchNotifier{
		inner:   inner,
		window:  window,
		maxSize: maxSize,
	}, nil
}

// Send adds msg to the current batch and flushes if maxSize is reached.
func (b *BatchNotifier) Send(ctx context.Context, msg Message) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.batch = append(b.batch, msg)

	if b.timer == nil {
		b.timer = time.AfterFunc(b.window, func() {
			b.mu.Lock()
			defer b.mu.Unlock()
			_ = b.flushLocked(ctx)
		})
	}

	if len(b.batch) >= b.maxSize {
		if b.timer != nil {
			b.timer.Stop()
			b.timer = nil
		}
		return b.flushLocked(ctx)
	}
	return nil
}

// Flush sends all pending messages immediately.
func (b *BatchNotifier) Flush(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.timer != nil {
		b.timer.Stop()
		b.timer = nil
	}
	return b.flushLocked(ctx)
}

func (b *BatchNotifier) flushLocked(ctx context.Context) error {
	if len(b.batch) == 0 {
		return nil
	}
	// Send the first message as representative; inner notifier receives it.
	summary := b.batch[0]
	summary.Summary = formatBatchSummary(b.batch)
	b.batch = nil
	return b.inner.Send(ctx, summary)
}

func formatBatchSummary(msgs []Message) string {
	if len(msgs) == 1 {
		return msgs[0].Summary
	}
	return fmt.Sprintf("Batch of %d alerts; latest: %s", len(msgs), msgs[len(msgs)-1].Summary)
}
