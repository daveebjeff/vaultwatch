package notify

import "errors"

// ErrZeroDuration is returned when a duration argument is zero or negative.
var ErrZeroDuration = errors.New("notify: duration must be greater than zero")
