package notify

import (
	"context"
	"fmt"
)

// SequenceNotifier sends a message through a chain of notifiers in order,
// stopping at the first failure. Unlike MultiNotifier, which attempts all
// notifiers regardless of errors, SequenceNotifier is strict: a failure in
// any step halts the chain and returns the error immediately.
type SequenceNotifier struct {
	steps []Notifier
}

// NewSequenceNotifier creates a SequenceNotifier that will deliver messages
// to each notifier in the provided order. At least one notifier is required,
// and no notifier may be nil.
func NewSequenceNotifier(steps ...Notifier) (*SequenceNotifier, error) {
	if len(steps) == 0 {
		return nil, fmt.Errorf("sequence: at least one notifier is required")
	}
	for i, s := range steps {
		if s == nil {
			return nil, fmt.Errorf("sequence: notifier at index %d is nil", i)
		}
	}
	return &SequenceNotifier{steps: steps}, nil
}

// Send delivers msg to each notifier in sequence. If any notifier returns an
// error, Send stops immediately and returns that error with the step index
// included in the message.
func (s *SequenceNotifier) Send(ctx context.Context, msg Message) error {
	for i, n := range s.steps {
		if err := n.Send(ctx, msg); err != nil {
			return fmt.Errorf("sequence: step %d failed: %w", i, err)
		}
	}
	return nil
}

// Len returns the number of steps in the sequence.
func (s *SequenceNotifier) Len() int {
	return len(s.steps)
}
