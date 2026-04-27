package notify

import "errors"

// ErrQuotaExceeded is returned by QuotaNotifier when the maximum number of
// notifications for the current window has been reached.
var ErrQuotaExceeded = errors.New("quota: notification quota exceeded for current window")
