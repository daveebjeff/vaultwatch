package notify

import (
	"fmt"
	"sync"
	"time"
)

// BufferNotifier accumulates messages and flushes them as a batch after a
// configurable window or when the buffer reaches a maximum size.
type BufferNotifier struct {
	mu       sync.Mutex
	inner    Notifier
	window   time.Duration
	maxSize  int
	buf      []Message
	timer    *time.Timer
}

// NewBufferNotifier returns a BufferNotifier that wraps inner.
// window is the maximum time to hold messages before flushing.
// maxSize is the maximum number of messages before an immediate flush (0 = unlimited).
func NewBufferNotifier(inner Notifier, window time.Duration, maxSize int) (*BufferNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("buffer: inner notifier must not be nil")
	}
	if window <= 0 {
		return nil, fmt.Errorf("buffer: window must be positive")
	}
	return &BufferNotifier{
		inner:   inner,
		window:  window,
		maxSize: maxSize,
	}, nil
}

// Send adds msg to the buffer. If the buffer reaches maxSize it flushes
// immediately; otherwise a timer ensures the buffer is flushed after window.
func (b *BufferNotifier) Send(msg Message) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.buf = append(b.buf, msg)

	if b.maxSize > 0 && len(b.buf) >= b.maxSize {
		if b.timer != nil {
			b.timer.Stop()
			b.timer = nil
		}
		return b.flush()
	}

	if b.timer == nil {
		b.timer = time.AfterFunc(b.window, func() {
			b.mu.Lock()
			defer b.mu.Unlock()
			_ = b.flush()
		})
	}
	return nil
}

// Flush forces an immediate send of all buffered messages.
func (b *BufferNotifier) Flush() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.flush()
}

// flush must be called with b.mu held.
func (b *BufferNotifier) flush() error {
	if len(b.buf) == 0 {
		return nil
	}
	if b.timer != nil {
		b.timer.Stop()
		b.timer = nil
	}
	var firstErr error
	for _, m := range b.buf {
		if err := b.inner.Send(m); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	b.buf = b.buf[:0]
	return firstErr
}
