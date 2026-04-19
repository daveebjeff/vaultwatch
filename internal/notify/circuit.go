package notify

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// CircuitState represents the state of the circuit breaker.
type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

// CircuitNotifier wraps a Notifier with a circuit breaker that opens after
// consecutive failures and resets after a cooldown period.
type CircuitNotifier struct {
	mu          sync.Mutex
	inner       Notifier
	maxFailures int
	resetAfter  time.Duration
	failures    int
	state       CircuitState
	openedAt    time.Time
}

// NewCircuitNotifier returns a CircuitNotifier that opens after maxFailures
// consecutive errors and attempts recovery after resetAfter duration.
func NewCircuitNotifier(inner Notifier, maxFailures int, resetAfter time.Duration) (*CircuitNotifier, error) {
	if inner == nil {
		return nil, errors.New("circuit: inner notifier must not be nil")
	}
	if maxFailures <= 0 {
		return nil, errors.New("circuit: maxFailures must be greater than zero")
	}
	if resetAfter <= 0 {
		return nil, errors.New("circuit: resetAfter must be greater than zero")
	}
	return &CircuitNotifier{
		inner:       inner,
		maxFailures: maxFailures,
		resetAfter:  resetAfter,
		state:       CircuitClosed,
	}, nil
}

// Send forwards the message if the circuit is closed or half-open.
func (c *CircuitNotifier) Send(msg Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch c.state {
	case CircuitOpen:
		if time.Since(c.openedAt) >= c.resetAfter {
			c.state = CircuitHalfOpen
		} else {
			return fmt.Errorf("circuit: open, skipping notification for %s", msg.Path)
		}
	case CircuitClosed, CircuitHalfOpen:
		// proceed
	}

	err := c.inner.Send(msg)
	if err != nil {
		c.failures++
		if c.failures >= c.maxFailures {
			c.state = CircuitOpen
			c.openedAt = time.Now()
		}
		return err
	}
	c.failures = 0
	c.state = CircuitClosed
	return nil
}

// State returns the current circuit state.
func (c *CircuitNotifier) State() CircuitState {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state
}
