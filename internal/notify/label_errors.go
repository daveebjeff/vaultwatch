package notify

import "errors"

// Sentinel errors for LabelNotifier construction.
var (
	ErrLabelNilInner  = errors.New("label: inner notifier must not be nil")
	ErrLabelNoLabels  = errors.New("label: at least one label is required")
)
