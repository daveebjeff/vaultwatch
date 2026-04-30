package notify

import (
	"errors"
	"sync"
)

// SieveNotifier routes each message to the first notifier whose predicate
// returns true. If no predicate matches, the message is silently dropped
// (no error). Predicates are evaluated in the order they were added.
type SieveNotifier struct {
	mu      sync.RWMutex
	routes  []sieveRoute
}

type sieveRoute struct {
	predicate func(Message) bool
	notifier  Notifier
}

// NewSieveNotifier returns an empty SieveNotifier. Use Add to register
// predicate/notifier pairs before calling Send.
func NewSieveNotifier() *SieveNotifier {
	return &SieveNotifier{}
}

// Add registers a predicate and the notifier that should receive messages
// for which the predicate returns true. Returns an error if either argument
// is nil.
func (s *SieveNotifier) Add(predicate func(Message) bool, n Notifier) error {
	if predicate == nil {
		return errors.New("sieve: predicate must not be nil")
	}
	if n == nil {
		return errors.New("sieve: notifier must not be nil")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.routes = append(s.routes, sieveRoute{predicate: predicate, notifier: n})
	return nil
}

// Send evaluates each registered predicate in order and forwards msg to the
// first matching notifier. If no route matches, Send returns nil.
func (s *SieveNotifier) Send(msg Message) error {
	s.mu.RLock()
	routes := make([]sieveRoute, len(s.routes))
	copy(routes, s.routes)
	s.mu.RUnlock()

	for _, r := range routes {
		if r.predicate(msg) {
			return r.notifier.Send(msg)
		}
	}
	return nil
}
