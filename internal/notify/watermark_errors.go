package notify

import "errors"

// ErrWatermarkNilInner is returned when a nil inner notifier is provided.
var ErrWatermarkNilInner = errors.New("watermark: inner notifier must not be nil")

// ErrWatermarkInvalidDuration is returned when the watermark duration is not positive.
var ErrWatermarkInvalidDuration = errors.New("watermark: duration must be positive")
