package notify

import (
	"strings"
	"unicode"
)

// NormalizeNotifier sanitizes outbound message bodies before forwarding.
// It trims leading/trailing whitespace, collapses internal runs of whitespace,
// and optionally converts the body to a canonical case.
type NormalizeNotifier struct {
	inner     Notifier
	lowerCase bool
}

// NormalizeOption configures a NormalizeNotifier.
type NormalizeOption func(*NormalizeNotifier)

// WithLowerCase converts the message body to lower-case during normalization.
func WithLowerCase() NormalizeOption {
	return func(n *NormalizeNotifier) {
		n.lowerCase = true
	}
}

// NewNormalizeNotifier returns a Notifier that normalizes message bodies
// before delegating to inner. Returns ErrNilInner if inner is nil.
func NewNormalizeNotifier(inner Notifier, opts ...NormalizeOption) (*NormalizeNotifier, error) {
	if inner == nil {
		return nil, ErrNilInner
	}
	n := &NormalizeNotifier{inner: inner}
	for _, o := range opts {
		o(n)
	}
	return n, nil
}

// Send normalizes msg.Body and forwards the result to the inner Notifier.
func (n *NormalizeNotifier) Send(msg Message) error {
	msg.Body = normalize(msg.Body, n.lowerCase)
	return n.inner.Send(msg)
}

// normalize trims, collapses whitespace, and optionally lower-cases s.
func normalize(s string, lower bool) string {
	s = strings.TrimSpace(s)
	var b strings.Builder
	inSpace := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !inSpace {
				b.WriteRune(' ')
				inSpace = true
			}
			continue
		}
		inSpace = false
		b.WriteRune(r)
	}
	result := b.String()
	if lower {
		return strings.ToLower(result)
	}
	return result
}
