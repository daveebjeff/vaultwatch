package notify

import "errors"

var (
	errNilSecondary = errors.New("notify: secondary notifier must not be nil")
	errZeroDuration = errors.New("notify: hedge delay must be greater than zero")
)
