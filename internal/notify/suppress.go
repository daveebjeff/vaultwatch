package notify

import (
	"fmt"
	"sync"
	"time"
)

// SuppressNotifier silences alerts for a specific secret path for a fixed
// duration. Useful when an operator has acknowledged an issue and wants to
// stop receiving repeated notifications temporarily.
type SuppressNotifier struct {
	inner      Notifier
	mu         sync.Mutex
	suppressed map[string]time.Time
	ttl        time.Duration
}

// NewSuppressNotifier wraps inner and suppresses repeated alerts for the same
// path for ttl duration after Suppress is called.
func NewSuppressNotifier(inner Notifier, ttl time.Duration) (*SuppressNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("suppress: inner notifier must not be nil")
	}
	if ttl <= 0 {
		return nil, fmt.Errorf("suppress: ttl must be positive")
	}
	return &SuppressNotifier{
		inner:      inner,
		suppressed: make(map[string]time.Time),
		ttl:        ttl,
	}, nil
}

// Suppress silences alerts for the given path for the configured TTL.
func (s *SuppressNotifier) Suppress(path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.suppressed[path] = time.Now().Add(s.ttl)
}

// Unsuppress lifts suppression for the given path immediately.
func (s *SuppressNotifier) Unsuppress(path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.suppressed, path)
}

// Send forwards the message to the inner notifier unless the path is currently
// suppressed.
func (s *SuppressNotifier) Send(msg Message) error {
	s.mu.Lock()
	expiry, ok := s.suppressed[msg.Path]
	if ok && time.Now().Before(expiry) {
		s.mu.Unlock()
		return nil
	}
	if ok {
		delete(s.suppressed, msg.Path)
	}
	s.mu.Unlock()
	return s.inner.Send(msg)
}
