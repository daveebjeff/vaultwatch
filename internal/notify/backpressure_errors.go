package notify

import "errors"

// Sentinel errors for BackpressureNotifier.
var (
	ErrBackpressureNilInner     = errors.New("backpressure: inner notifier must not be nil")
	ErrBackpressureZeroCapacity = errors.New("backpressure: capacity must be greater than zero")
	ErrBackpressureQueueFull    = errors.New("backpressure: queue is full")
)
