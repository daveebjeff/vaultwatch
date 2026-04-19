package notify

import (
	"context"
	"log"
)

// ShadowNotifier sends to a primary notifier and mirrors to a shadow notifier
// for comparison or testing purposes. Errors from the shadow are logged but
// never returned to the caller.
type ShadowNotifier struct {
	primary Notifier
	shadow  Notifier
	logger  *log.Logger
}

// NewShadowNotifier creates a ShadowNotifier that forwards all messages to
// primary and mirrors them to shadow. Shadow errors are suppressed.
func NewShadowNotifier(primary, shadow Notifier, logger *log.Logger) (*ShadowNotifier, error) {
	if primary == nil {
		return nil, ErrNilNotifier
	}
	if shadow == nil {
		return nil, ErrNilNotifier
	}
	if logger == nil {
		logger = log.Default()
	}
	return &ShadowNotifier{
		primary: primary,
		shadow:  shadow,
		logger:  logger,
	}, nil
}

// Send delivers the message to the primary notifier and asynchronously mirrors
// it to the shadow notifier. Only primary errors are returned.
func (s *ShadowNotifier) Send(ctx context.Context, msg Message) error {
	go func() {
		if err := s.shadow.Send(ctx, msg); err != nil {
			s.logger.Printf("shadow notifier error for path %s: %v", msg.Path, err)
		}
	}()
	return s.primary.Send(ctx, msg)
}
