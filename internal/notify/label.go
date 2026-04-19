package notify

import "fmt"

// LabelNotifier wraps a Notifier and attaches static key/value labels
// to every outgoing message's summary for easier downstream filtering.
type LabelNotifier struct {
	inner  Notifier
	labels map[string]string
}

// NewLabelNotifier returns a LabelNotifier that prepends labels to each
// message summary before forwarding to inner.
func NewLabelNotifier(inner Notifier, labels map[string]string) (*LabelNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("label: inner notifier must not be nil")
	}
	if len(labels) == 0 {
		return nil, fmt.Errorf("label: at least one label is required")
	}
	copy := make(map[string]string, len(labels))
	for k, v := range labels {
		copy[k] = v
	}
	return &LabelNotifier{inner: inner, labels: copy}, nil
}

// Send attaches labels to the message summary and forwards to the inner notifier.
func (l *LabelNotifier) Send(msg Message) error {
	tagged := msg
	prefix := ""
	for k, v := range l.labels {
		prefix += fmt.Sprintf("[%s=%s] ", k, v)
	}
	tagged.Summary = prefix + msg.Summary
	return l.inner.Send(tagged)
}
