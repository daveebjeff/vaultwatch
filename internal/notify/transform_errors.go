package notify

import "errors"

// Sentinel errors returned by NewTransformNotifier.
var (
	// ErrTransformNilInner is returned when a nil inner notifier is provided.
	ErrTransformNilInner = errors.New("transform: inner notifier must not be nil")

	// ErrTransformNilFn is returned when a nil transform function is provided.
	ErrTransformNilFn = errors.New("transform: transform function must not be nil")
)
