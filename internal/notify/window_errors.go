package notify

import "errors"

// ErrWindowLimitExceeded is returned by WindowNotifier when the sliding-window
// ceiling has been reached.  Callers may check for this with errors.Is.
var ErrWindowLimitExceeded = errors.New("window: rate limit exceeded")
