package notify

import "errors"

var (
	errNilInner    = errors.New("notify: inner notifier must not be nil")
	errInvalidRate = errors.New("notify: sampling rate must be between 0.0 and 1.0")
)
