package notify

import "errors"

// ErrTraceIDNilInner is returned when NewTraceIDNotifier receives a nil inner notifier.
var ErrTraceIDNilInner = errors.New("traceid: inner notifier must not be nil")
