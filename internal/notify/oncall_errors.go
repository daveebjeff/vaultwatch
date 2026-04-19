package notify

import "errors"

// ErrNoOnCallRotation is returned when Send is called but no rotation window
// matches the current time.
var ErrNoOnCallRotation = errors.New("oncall: no active rotation for current time")
