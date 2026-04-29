package notify

import (
	"context"
	"sync"
)

// OnceNotifier forwards a message to the inner notifier exactly once per
// unique secret path. Subsequent sends for the same path are silently
// dropped. Call Reset to clear the seen set and allow re-delivery.
//
// This is useful when you want a single fire-and-forget alert per secret,
// regardless of how many monitor cycles observe the same condition.
type OnceNotifier struct {
	inner Notifier
	mu   sync.Mutex
	seen map[string]struct{}
}

// NewOnceNotifier returns a OnceNotifier wrapping inner.
// It returns an error if inner is nil.
func NewOnceNotifier(inner Notifier) (*OnceNotifier, error) {
	if inner == nil {
		return nil, ErrNilInner
	}
	return &OnceNotifier{
		inner: inner,
		seen:  make(map[string]struct{}),
	}, nil
}

// Send forwards msg to the inner notifier only if this path has not been
// seen before. Returns nil (without forwarding) when the path is already
// recorded.
func (n *OnceNotifier) Send(ctx context.Context, msg Message) error {
	n.mu.Lock()
	_, already := n.seen[msg.Path]
	if !already {
		n.seen[msg.Path] = struct{}{}
	}
	n.mu.Unlock()

	if already {
		return nil
	}
	return n.inner.Send(ctx, msg)
}

// Reset clears the set of seen paths so that every path may be delivered
// once more.
func (n *OnceNotifier) Reset() {
	n.mu.Lock()
	n.seen = make(map[string]struct{})
	n.mu.Unlock()
}

// Seen reports whether the given path has already been delivered.
func (n *OnceNotifier) Seen(path string) bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	_, ok := n.seen[path]
	return ok
}
