package notify

import "errors"

// errInvalidMax is returned when maxItems is not a positive integer.
var errInvalidMax = errors.New("notify: maxItems must be greater than zero")

// errInvalidWindow is returned when maxAge is not a positive duration.
var errInvalidWindow = errors.New("notify: window/maxAge must be greater than zero")
