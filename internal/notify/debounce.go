package notify

import (
	"sync"
	"time"
)

// DebounceNotifier delays forwarding a notification until no new notification
// for the same secret path has arrived within the given wait window. If a new
// notification arrives before the timer fires, the timer resets. Only the
// most-recent message is forwarded when the timer finally fires.
type DebounceNotifier struct {
	inner  Notifier
	wait   time.Duration
	mu     sync.Mutex
	timers map[string]*time.Timer
	latest map[string]Message
}

// NewDebounceNotifier returns a DebounceNotifier that waits for wait duration
// of silence on a given secret path before forwarding to inner.
func NewDebounceNotifier(inner Notifier, wait time.Duration) (*DebounceNotifier, error) {
	if inner == nil {
		return nil, ErrNilInner
	}
	if wait <= 0 {
		return nil, ErrZeroWindow
	}
	return &DebounceNotifier{
		inner:  inner,
		wait:   wait,
		timers: make(map[string]*time.Timer),
		latest: make(map[string]Message),
	}, nil
}

// Send resets the debounce timer for msg.Path, storing msg as the latest
// message. When the timer expires the most-recent message is forwarded.
func (d *DebounceNotifier) Send(msg Message) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.latest[msg.Path] = msg

	if t, ok := d.timers[msg.Path]; ok {
		t.Reset(d.wait)
		return nil
	}

	d.timers[msg.Path] = time.AfterFunc(d.wait, func() {
		d.mu.Lock()
		m := d.latest[msg.Path]
		delete(d.timers, msg.Path)
		delete(d.latest, msg.Path)
		d.mu.Unlock()
		_ = d.inner.Send(m)
	})
	return nil
}
