package notify

import (
	"context"
	"errors"
)

// ConditionalNotifier forwards messages to inner only when the predicate returns true.
type ConditionalNotifier struct {
	inner     Notifier
	predicate func(Message) bool
}

// NewConditionalNotifier creates a ConditionalNotifier that forwards messages
// to inner when predicate returns true. Returns an error if inner is nil or
// predicate is nil.
func NewConditionalNotifier(inner Notifier, predicate func(Message) bool) (*ConditionalNotifier, error) {
	if inner == nil {
		return nil, errors.New("conditional: inner notifier must not be nil")
	}
	if predicate == nil {
		return nil, errors.New("conditional: predicate must not be nil")
	}
	return &ConditionalNotifier{inner: inner, predicate: predicate}, nil
}

// Send forwards msg to the inner notifier only if the predicate returns true.
func (c *ConditionalNotifier) Send(ctx context.Context, msg Message) error {
	if !c.predicate(msg) {
		return nil
	}
	return c.inner.Send(ctx, msg)
}
