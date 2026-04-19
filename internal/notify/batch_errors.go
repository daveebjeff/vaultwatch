package notify

import "errors"

// ErrZeroWindow is returned when a zero or negative window duration is provided.
var ErrZeroWindow = errors.New("notify: window duration must be greater than zero")
