package notify

import (
	"context"
	"sync"
	"time"
)

// StickyNotifier forwards a message and continues re-sending it at the given
// interval until the path is explicitly cleared. This is useful for alerts
// that should repeat until an operator acknowledges the condition.
type StickyNotifier struct {
	inner    Notifier
	interval time.Duration

	mu      sync.Mutex
	active  map[string]context.CancelFunc
	last    map[string]Message
}

// NewStickyNotifier creates a StickyNotifier that re-fires messages for active
// paths every interval until Clear is called for that path.
func NewStickyNotifier(inner Notifier, interval time.Duration) (*StickyNotifier, error) {
	if inner == nil {
		return nil, ErrNilInner
	}
	if interval <= 0 {
		return nil, ErrZeroWindow
	}
	return &StickyNotifier{
		inner:    inner,
		interval: interval,
		active:   make(map[string]context.CancelFunc),
		last:     make(map[string]Message),
	}, nil
}

// Send forwards the message immediately and starts a background ticker that
// re-sends the same message every interval. If the path is already sticky the
// ticker is restarted with the newest message.
func (s *StickyNotifier) Send(ctx context.Context, msg Message) error {
	if err := s.inner.Send(ctx, msg); err != nil {
		return err
	}

	s.mu.Lock()
	if cancel, ok := s.active[msg.Path]; ok {
		cancel()
	}
	tickCtx, cancel := context.WithCancel(context.Background())
	s.active[msg.Path] = cancel
	s.last[msg.Path] = msg
	s.mu.Unlock()

	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()
		for {
			select {
			case <-tickCtx.Done():
				return
			case <-ticker.C:
				s.mu.Lock()
				m := s.last[msg.Path]
				s.mu.Unlock()
				_ = s.inner.Send(tickCtx, m)
			}
		}
	}()

	return nil
}

// Clear stops the repeating ticker for the given path.
func (s *StickyNotifier) Clear(path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if cancel, ok := s.active[path]; ok {
		cancel()
		delete(s.active, path)
		delete(s.last, path)
	}
}

// ActivePaths returns the set of paths currently being repeated.
func (s *StickyNotifier) ActivePaths() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	paths := make([]string, 0, len(s.active))
	for p := range s.active {
		paths = append(paths, p)
	}
	return paths
}
