package notify

import (
	"sync"
	"time"
)

// CoalesceNotifier groups rapid successive notifications for the same secret
// path into a single delivery. Within the coalesce window, only the most
// recent message is forwarded when the window closes. Each path has its own
// independent timer.
type CoalesceNotifier struct {
	inner    Notifier
	window   time.Duration
	mu       sync.Mutex
	pending  map[string]*coalesceEntry
}

type coalesceEntry struct {
	msg   Message
	timer *time.Timer
}

// NewCoalesceNotifier creates a CoalesceNotifier that waits window duration
// after the last message for a given path before forwarding it to inner.
// This ensures bursts of updates collapse into a single notification.
func NewCoalesceNotifier(inner Notifier, window time.Duration) (*CoalesceNotifier, error) {
	if inner == nil {
		return nil, ErrNilInner
	}
	if window <= 0 {
		return nil, ErrZeroWindow
	}
	return &CoalesceNotifier{
		inner:   inner,
		window:  window,
		pending: make(map[string]*coalesceEntry),
	}, nil
}

// Send schedules msg for delivery after the coalesce window. If a pending
// message for the same path already exists its timer is reset and the message
// is replaced with the newer one.
func (c *CoalesceNotifier) Send(msg Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, ok := c.pending[msg.Path]; ok {
		entry.timer.Stop()
		entry.msg = msg
		entry.timer = time.AfterFunc(c.window, c.flush(msg.Path))
		return nil
	}

	c.pending[msg.Path] = &coalesceEntry{
		msg:   msg,
		timer: time.AfterFunc(c.window, c.flush(msg.Path)),
	}
	return nil
}

func (c *CoalesceNotifier) flush(path string) func() {
	return func() {
		c.mu.Lock()
		entry, ok := c.pending[path]
		if ok {
			delete(c.pending, path)
		}
		c.mu.Unlock()
		if ok {
			_ = c.inner.Send(entry.msg)
		}
	}
}
