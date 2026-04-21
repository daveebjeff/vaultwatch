package notify

import "fmt"

// TransformNotifier applies a user-supplied transform function to each
// Message before forwarding it to the inner Notifier. This allows
// callers to enrich, redact, or rewrite messages in a pipeline.
type TransformNotifier struct {
	inner Notifier
	fn    func(Message) Message
}

// NewTransformNotifier creates a TransformNotifier that applies fn to
// every message before passing it to inner.
// Returns an error if inner is nil or fn is nil.
func NewTransformNotifier(inner Notifier, fn func(Message) Message) (*TransformNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("transform: inner notifier must not be nil")
	}
	if fn == nil {
		return nil, fmt.Errorf("transform: transform function must not be nil")
	}
	return &TransformNotifier{inner: inner, fn: fn}, nil
}

// Send applies the transform function to msg and forwards the result.
func (t *TransformNotifier) Send(msg Message) error {
	return t.inner.Send(t.fn(msg))
}
