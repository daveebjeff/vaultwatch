package notify

import (
	"context"
	"fmt"
)

// FallbackNotifier tries a primary notifier and, on failure, falls back to a
// secondary notifier. Unlike EscalationNotifier it is synchronous and does not
// use a timeout — it simply forwards the error from the secondary if both fail.
type FallbackNotifier struct {
	primary   Notifier
	secondary Notifier
}

// NewFallbackNotifier returns a FallbackNotifier that sends via primary and,
// if that fails, via secondary.
func NewFallbackNotifier(primary, secondary Notifier) (*FallbackNotifier, error) {
	if primary == nil {
		return nil, fmt.Errorf("fallback: primary notifier must not be nil")
	}
	if secondary == nil {
		return nil, fmt.Errorf("fallback: secondary notifier must not be nil")
	}
	return &FallbackNotifier{primary: primary, secondary: secondary}, nil
}

// Send attempts the primary notifier. On error it attempts the secondary and
// returns a combined error if the secondary also fails.
func (f *FallbackNotifier) Send(ctx context.Context, msg Message) error {
	if err := f.primary.Send(ctx, msg); err != nil {
		if err2 := f.secondary.Send(ctx, msg); err2 != nil {
			return fmt.Errorf("fallback: primary error: %w; secondary error: %v", err, err2)
		}
	}
	return nil
}
