package notify

import (
	"context"
	"errors"
	"fmt"
)

// PreSendNotifier runs a hook function before forwarding a message to the
// inner notifier. If the hook returns an error the message is not forwarded
// and the error is returned to the caller. This is useful for validation,
// enrichment, or conditional gating that needs to run synchronously before
// delivery.
type PreSendNotifier struct {
	inner  Notifier
	hook   func(ctx context.Context, msg Message) error
}

// NewPreSendNotifier creates a PreSendNotifier that calls hook before every
// Send. If hook returns a non-nil error the message is dropped and the error
// is returned. inner and hook must not be nil.
func NewPreSendNotifier(inner Notifier, hook func(ctx context.Context, msg Message) error) (*PreSendNotifier, error) {
	if inner == nil {
		return nil, errors.New("presend: inner notifier must not be nil")
	}
	if hook == nil {
		return nil, errors.New("presend: hook must not be nil")
	}
	return &PreSendNotifier{inner: inner, hook: hook}, nil
}

// Send calls the hook first. If the hook succeeds the message is forwarded
// to the inner notifier. If the hook fails the error is wrapped and returned
// without calling the inner notifier.
func (p *PreSendNotifier) Send(ctx context.Context, msg Message) error {
	if err := p.hook(ctx, msg); err != nil {
		return fmt.Errorf("presend hook rejected message: %w", err)
	}
	return p.inner.Send(ctx, msg)
}
