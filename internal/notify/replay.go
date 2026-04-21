package notify

import (
	"sync"
	"time"
)

// ReplayNotifier stores recent messages and can replay them to a new notifier.
// Useful for onboarding new downstream targets that need historical context.
type ReplayNotifier struct {
	inner    Notifier
	mu       sync.Mutex
	buffer   []Message
	maxAge   time.Duration
	maxItems int
}

// NewReplayNotifier wraps inner and retains up to maxItems messages younger
// than maxAge. Call Replay to forward retained messages to any Notifier.
func NewReplayNotifier(inner Notifier, maxItems int, maxAge time.Duration) (*ReplayNotifier, error) {
	if inner == nil {
		return nil, ErrNilInner
	}
	if maxItems <= 0 {
		return nil, errInvalidMax
	}
	if maxAge <= 0 {
		return nil, errInvalidWindow
	}
	return &ReplayNotifier{
		inner:    inner,
		maxItems: maxItems,
		maxAge:   maxAge,
	}, nil
}

// Send forwards the message to the inner notifier and retains it for replay.
func (r *ReplayNotifier) Send(msg Message) error {
	r.mu.Lock()
	r.prune()
	r.buffer = append(r.buffer, msg)
	if len(r.buffer) > r.maxItems {
		r.buffer = r.buffer[len(r.buffer)-r.maxItems:]
	}
	r.mu.Unlock()
	return r.inner.Send(msg)
}

// Replay forwards all retained messages to dst in chronological order.
func (r *ReplayNotifier) Replay(dst Notifier) error {
	r.mu.Lock()
	r.prune()
	snap := make([]Message, len(r.buffer))
	copy(snap, r.buffer)
	r.mu.Unlock()

	for _, msg := range snap {
		if err := dst.Send(msg); err != nil {
			return err
		}
	}
	return nil
}

// Len returns the number of currently retained messages.
func (r *ReplayNotifier) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prune()
	return len(r.buffer)
}

// prune removes messages older than maxAge. Caller must hold r.mu.
func (r *ReplayNotifier) prune() {
	cutoff := time.Now().Add(-r.maxAge)
	i := 0
	for i < len(r.buffer) && r.buffer[i].ExpiresAt.Before(cutoff) {
		i++
	}
	r.buffer = r.buffer[i:]
}
