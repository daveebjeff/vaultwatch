package notify

import (
	"fmt"
	"unicode/utf8"
)

// TruncateNotifier wraps a Notifier and truncates the Message body to a
// maximum number of runes before forwarding. This is useful when downstream
// channels (e.g. SMS, PagerDuty) impose strict character limits.
type TruncateNotifier struct {
	inner   Notifier
	maxRune int
	suffix  string
}

// NewTruncateNotifier returns a TruncateNotifier that truncates Message.Body
// to maxRunes characters. If the body is truncated, suffix is appended (e.g.
// "…"). Returns an error if inner is nil, maxRunes < 1, or suffix alone
// exceeds maxRunes.
func NewTruncateNotifier(inner Notifier, maxRunes int, suffix string) (*TruncateNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("truncate: inner notifier must not be nil")
	}
	if maxRunes < 1 {
		return nil, fmt.Errorf("truncate: maxRunes must be >= 1, got %d", maxRunes)
	}
	if utf8.RuneCountInString(suffix) >= maxRunes {
		return nil, fmt.Errorf("truncate: suffix length must be less than maxRunes")
	}
	return &TruncateNotifier{inner: inner, maxRune: maxRunes, suffix: suffix}, nil
}

// Send truncates msg.Body if it exceeds the configured limit, then forwards
// the (possibly modified) message to the inner notifier.
func (t *TruncateNotifier) Send(msg Message) error {
	if utf8.RuneCountInString(msg.Body) > t.maxRune {
		runes := []rune(msg.Body)
		cutAt := t.maxRune - utf8.RuneCountInString(t.suffix)
		msg.Body = string(runes[:cutAt]) + t.suffix
	}
	return t.inner.Send(msg)
}
