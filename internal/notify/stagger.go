package notify

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// StaggerNotifier sends notifications to a list of inner notifiers with a
// fixed delay between each delivery. This spreads load when many notifiers
// are configured and avoids thundering-herd on downstream services.
type StaggerNotifier struct {
	mu      sync.Mutex
	inners  []Notifier
	delay   time.Duration
}

// NewStaggerNotifier creates a StaggerNotifier that waits delay between
// each successive notifier call. At least one notifier must be provided
// and delay must be positive.
func NewStaggerNotifier(delay time.Duration, notifiers ...Notifier) (*StaggerNotifier, error) {
	if len(notifiers) == 0 {
		return nil, fmt.Errorf("stagger: at least one notifier required")
	}
	for i, n := range notifiers {
		if n == nil {
			return nil, fmt.Errorf("stagger: notifier at index %d is nil", i)
		}
	}
	if delay <= 0 {
		return nil, fmt.Errorf("stagger: delay must be positive, got %s", delay)
	}
	return &StaggerNotifier{
		inners: notifiers,
		delay:  delay,
	}, nil
}

// Send delivers msg to each inner notifier in order, sleeping delay between
// each call. If the context is cancelled mid-flight the send is aborted and
// the context error is returned. The first notifier error encountered is
// returned; subsequent notifiers are still attempted.
func (s *StaggerNotifier) Send(ctx context.Context, msg Message) error {
	s.mu.Lock()
	inners := make([]Notifier, len(s.inners))
	copy(inners, s.inners)
	delay := s.delay
	s.mu.Unlock()

	var firstErr error
	for i, n := range inners {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := n.Send(ctx, msg); err != nil && firstErr == nil {
			firstErr = err
		}
		if i < len(inners)-1 {
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	return firstErr
}

// Add appends a notifier to the stagger chain at runtime.
func (s *StaggerNotifier) Add(n Notifier) error {
	if n == nil {
		return fmt.Errorf("stagger: cannot add nil notifier")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.inners = append(s.inners, n)
	return nil
}
