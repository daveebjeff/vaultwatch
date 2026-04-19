package notify

import (
	"fmt"
	"sync"
	"time"
)

// OnCallNotifier routes alerts to the currently on-call notifier based on a
// rotation schedule. Rotations are evaluated at send time.
type OnCallNotifier struct {
	mu        sync.RWMutex
	rotations []OnCallRotation
	now       func() time.Time
}

// OnCallRotation pairs a time window with a notifier responsible for that slot.
type OnCallRotation struct {
	Name     string
	Start    time.Time // UTC
	End      time.Time // UTC
	Notifier Notifier
}

// NewOnCallNotifier creates an OnCallNotifier from the provided rotations.
// Returns an error if any rotation has a nil notifier or invalid window.
func NewOnCallNotifier(rotations []OnCallRotation) (*OnCallNotifier, error) {
	if len(rotations) == 0 {
		return nil, fmt.Errorf("oncall: at least one rotation is required")
	}
	for i, r := range rotations {
		if r.Notifier == nil {
			return nil, fmt.Errorf("oncall: rotation %d has nil notifier", i)
		}
		if !r.End.After(r.Start) {
			return nil, fmt.Errorf("oncall: rotation %d end must be after start", i)
		}
	}
	return &OnCallNotifier{rotations: rotations, now: time.Now}, nil
}

// Send delivers msg to the first rotation whose window contains the current time.
// If no rotation is active, ErrNoOnCallRotation is returned.
func (o *OnCallNotifier) Send(msg Message) error {
	o.mu.RLock()
	defer o.mu.RUnlock()
	now := o.now().UTC()
	for _, r := range o.rotations {
		if !now.Before(r.Start) && now.Before(r.End) {
			return r.Notifier.Send(msg)
		}
	}
	return ErrNoOnCallRotation
}

// AddRotation appends a new rotation at runtime.
func (o *OnCallNotifier) AddRotation(r OnCallRotation) error {
	if r.Notifier == nil {
		return fmt.Errorf("oncall: rotation has nil notifier")
	}
	if !r.End.After(r.Start) {
		return fmt.Errorf("oncall: end must be after start")
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	o.rotations = append(o.rotations, r)
	return nil
}
