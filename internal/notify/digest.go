package notify

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// DigestNotifier batches messages over a time window and sends a single
// summarised notification when the window closes or the max size is reached.
type DigestNotifier struct {
	mu       sync.Mutex
	inner    Notifier
	window   time.Duration
	maxSize  int
	pending  []Message
	timer    *time.Timer
}

// NewDigestNotifier creates a DigestNotifier that flushes after window duration
// or when maxSize messages have accumulated, whichever comes first.
func NewDigestNotifier(inner Notifier, window time.Duration, maxSize int) (*DigestNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("digest: inner notifier must not be nil")
	}
	if window <= 0 {
		return nil, fmt.Errorf("digest: window must be positive")
	}
	if maxSize <= 0 {
		maxSize = 50
	}
	return &DigestNotifier{
		inner:   inner,
		window:  window,
		maxSize: maxSize,
	}, nil
}

// Send enqueues the message and flushes if the batch is full.
func (d *DigestNotifier) Send(msg Message) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if len(d.pending) == 0 {
		d.timer = time.AfterFunc(d.window, func() {
			d.mu.Lock()
			defer d.mu.Unlock()
			_ = d.flush()
		})
	}

	d.pending = append(d.pending, msg)

	if len(d.pending) >= d.maxSize {
		if d.timer != nil {
			d.timer.Stop()
		}
		return d.flush()
	}
	return nil
}

// Flush forces immediate delivery of all pending messages.
func (d *DigestNotifier) Flush() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.timer != nil {
		d.timer.Stop()
	}
	return d.flush()
}

// flush must be called with d.mu held.
func (d *DigestNotifier) flush() error {
	if len(d.pending) == 0 {
		return nil
	}
	lines := make([]string, 0, len(d.pending))
	worst := StatusOK
	for _, m := range d.pending {
		lines = append(lines, fmt.Sprintf("[%s] %s — %s", m.Status, m.Path, m.Detail))
		if m.Status > worst {
			worst = m.Status
		}
	}
	summary := Message{
		Path:   "digest",
		Status: worst,
		Detail: fmt.Sprintf("%d secrets:\n%s", len(lines), strings.Join(lines, "\n")),
	}
	d.pending = d.pending[:0]
	return d.inner.Send(summary)
}
