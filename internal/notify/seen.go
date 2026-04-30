package notify

import (
	"context"
	"sync"
	"time"
)

// SeenNotifier suppresses messages for paths that have already been seen
// within a rolling time window. Unlike DedupNotifier, it tracks the first
// occurrence time and only re-forwards once the window has fully elapsed.
type SeenNotifier struct {
	inner   Notifier
	window  time.Duration
	mu      sync.Mutex
	firstAt map[string]time.Time
}

// NewSeenNotifier creates a SeenNotifier that wraps inner and suppresses
// repeat messages for the same path within window. Returns an error if
// inner is nil or window is zero.
func NewSeenNotifier(inner Notifier, window time.Duration) (*SeenNotifier, error) {
	if inner == nil {
		return nil, ErrNilInner
	}
	if window <= 0 {
		return nil, ErrZeroWindow
	}
	return &SeenNotifier{
		inner:   inner,
		window:  window,
		firstAt: make(map[string]time.Time),
	}, nil
}

// Send forwards the message to the inner notifier only if the path has not
// been seen within the configured window. The first occurrence resets the
// window clock for that path.
func (s *SeenNotifier) Send(ctx context.Context, msg Message) error {
	s.mu.Lock()
	now := time.Now()
	first, exists := s.firstAt[msg.Path]
	if !exists || now.Sub(first) >= s.window {
		s.firstAt[msg.Path] = now
		s.mu.Unlock()
		return s.inner.Send(ctx, msg)
	}
	s.mu.Unlock()
	return nil
}

// Forget removes the path from the seen map, allowing the next message for
// that path to be forwarded regardless of the window.
func (s *SeenNotifier) Forget(path string) {
	s.mu.Lock()
	delete(s.firstAt, path)
	s.mu.Unlock()
}

// Reset clears all tracked paths.
func (s *SeenNotifier) Reset() {
	s.mu.Lock()
	s.firstAt = make(map[string]time.Time)
	s.mu.Unlock()
}
