package notify

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// RollupNotifier collects messages over a window and sends a single
// summary notification listing all affected secret paths.
type RollupNotifier struct {
	mu       sync.Mutex
	inner    Notifier
	window   time.Duration
	msgs     []Message
	timer    *time.Timer
	maxSize  int
}

// NewRollupNotifier creates a RollupNotifier that flushes after window
// duration or when maxSize messages have accumulated.
func NewRollupNotifier(inner Notifier, window time.Duration, maxSize int) (*RollupNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("rollup: inner notifier must not be nil")
	}
	if window <= 0 {
		return nil, fmt.Errorf("rollup: window must be positive")
	}
	if maxSize <= 0 {
		maxSize = 50
	}
	return &RollupNotifier{inner: inner, window: window, maxSize: maxSize}, nil
}

// Send buffers the message and triggers a flush when the window or size limit is reached.
func (r *RollupNotifier) Send(msg Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.msgs = append(r.msgs, msg)

	if r.timer == nil {
		r.timer = time.AfterFunc(r.window, func() {
			r.mu.Lock()
			defer r.mu.Unlock()
			r.flush()
		})
	}

	if len(r.msgs) >= r.maxSize {
		if r.timer != nil {
			r.timer.Stop()
			r.timer = nil
		}
		r.flush()
	}
	return nil
}

// Flush forces an immediate send of buffered messages.
func (r *RollupNotifier) Flush() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.timer != nil {
		r.timer.Stop()
		r.timer = nil
	}
	return r.flush()
}

func (r *RollupNotifier) flush() error {
	if len(r.msgs) == 0 {
		return nil
	}
	defer func() { r.msgs = nil }()

	var lines []string
	worst := StatusOK
	for _, m := range r.msgs {
		lines = append(lines, fmt.Sprintf("  [%s] %s", m.Status, m.SecretPath))
		if m.Status > worst {
			worst = m.Status
		}
	}
	summary := Message{
		SecretPath: "(rollup)",
		Status:     worst,
		Summary:    fmt.Sprintf("%d secret(s) need attention:\n%s", len(r.msgs), strings.Join(lines, "\n")),
	}
	return r.inner.Send(summary)
}
