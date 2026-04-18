package notify

import (
	"fmt"
	"sync"
)

// DedupNotifier wraps a Notifier and suppresses consecutive duplicate
// messages for the same path and status, forwarding only when something
// changes.
type DedupNotifier struct {
	mu    sync.Mutex
	inner Notifier
	last  map[string]Status
}

// NewDedupNotifier returns a DedupNotifier wrapping inner.
func NewDedupNotifier(inner Notifier) (*DedupNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("dedup: inner notifier must not be nil")
	}
	return &DedupNotifier{
		inner: inner,
		last:  make(map[string]Status),
	}, nil
}

// Send forwards msg to the inner notifier only if the status for msg.Path
// has changed since the last forwarded message.
func (d *DedupNotifier) Send(msg Message) error {
	d.mu.Lock()
	prev, seen := d.last[msg.Path]
	if seen && prev == msg.Status {
		d.mu.Unlock()
		return nil
	}
	d.last[msg.Path] = msg.Status
	d.mu.Unlock()
	return d.inner.Send(msg)
}
