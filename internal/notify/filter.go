package notify

import "strings"

// FilterNotifier wraps a Notifier and only forwards messages whose paths
// match at least one of the configured prefix patterns.
type FilterNotifier struct {
	inner    Notifier
	prefixes []string
}

// NewFilterNotifier returns a FilterNotifier that forwards to inner only when
// the message path matches one of the given prefixes. At least one prefix and
// a non-nil inner notifier are required.
func NewFilterNotifier(inner Notifier, prefixes []string) (*FilterNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("filter: inner notifier must not be nil")
	}
	if len(prefixes) == 0 {
		return nil, fmt.Errorf("filter: at least one prefix is required")
	}
	return &FilterNotifier{inner: inner, prefixes: prefixes}, nil
}

// Send forwards msg to the inner notifier if the message path matches any
// configured prefix. Returns nil without sending if no prefix matches.
func (f *FilterNotifier) Send(msg Message) error {
	for _, p := range f.prefixes {
		if strings.HasPrefix(msg.Path, p) {
			return f.inner.Send(msg)
		}
	}
	return nil
}
