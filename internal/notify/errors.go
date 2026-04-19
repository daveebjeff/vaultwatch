package notify

import "errors"

// Sentinel errors shared across notifier wrappers.
var (
	// ErrNilInner is returned when a wrapper is constructed with a nil inner Notifier.
	ErrNilInner = errors.New("notify: inner notifier must not be nil")

	// ErrZeroCooldown is returned when a cooldown or TTL duration is zero or negative.
	ErrZeroCooldown = errors.New("notify: cooldown/ttl duration must be greater than zero")

	// ErrZeroWindow is returned when a time window is zero or negative.
	ErrZeroWindow = errors.New("notify: window duration must be greater than zero")

	// ErrZeroMax is returned when a maximum count is zero or negative.
	ErrZeroMax = errors.New("notify: max count must be greater than zero")
)
