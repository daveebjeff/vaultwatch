package notify

import (
	"context"
	"regexp"
	"strings"
)

// RedactNotifier replaces sensitive patterns in a message body before
// forwarding to the inner notifier. Useful for scrubbing tokens, passwords,
// or other secrets that may appear in alert payloads.
type RedactNotifier struct {
	inner    Notifier
	patterns []*regexp.Regexp
	replacement string
}

// NewRedactNotifier returns a RedactNotifier that applies each compiled
// pattern to the message body, replacing matches with replacement.
// replacement defaults to "[REDACTED]" when empty.
func NewRedactNotifier(inner Notifier, patterns []*regexp.Regexp, replacement string) (*RedactNotifier, error) {
	if inner == nil {
		return nil, ErrNilInner
	}
	if len(patterns) == 0 {
		return nil, errRedactNoPatterns
	}
	if replacement == "" {
		replacement = "[REDACTED]"
	}
	return &RedactNotifier{
		inner:       inner,
		patterns:    patterns,
		replacement: replacement,
	}, nil
}

// Send redacts the message body and forwards the modified message.
func (r *RedactNotifier) Send(ctx context.Context, msg Message) error {
	body := msg.Body
	for _, re := range r.patterns {
		body = re.ReplaceAllString(body, r.replacement)
	}
	msg.Body = body
	return r.inner.Send(ctx, msg)
}

// MustCompilePatterns compiles a slice of raw regex strings and panics on
// the first invalid pattern. Intended for use in tests or init blocks.
func MustCompilePatterns(raw []string) []*regexp.Regexp {
	out := make([]*regexp.Regexp, len(raw))
	for i, s := range raw {
		out[i] = regexp.MustCompile(s)
	}
	return out
}

// CompilePatterns compiles a slice of raw regex strings, returning the first
// error encountered.
func CompilePatterns(raw []string) ([]*regexp.Regexp, error) {
	out := make([]*regexp.Regexp, 0, len(raw))
	for _, s := range raw {
		re, err := regexp.Compile(s)
		if err != nil {
			return nil, err
		}
		out = append(out, re)
	}
	return out, nil
}

// defaultRedactPatterns are common secret-like patterns applied when no
// explicit patterns are provided by callers using helper constructors.
var defaultRedactPatterns = MustCompilePatterns([]string{
	`(?i)(token|password|secret|key)=[^\s&]+`,
	`s\.[A-Za-z0-9]{24,}`, // Vault token format
})

// NewDefaultRedactNotifier wraps inner with the built-in secret patterns.
func NewDefaultRedactNotifier(inner Notifier) (*RedactNotifier, error) {
	return NewRedactNotifier(inner, defaultRedactPatterns, "")
}

// redactBody is a convenience used in tests.
func redactBody(patterns []*regexp.Regexp, replacement, body string) string {
	for _, re := range patterns {
		body = re.ReplaceAllString(body, replacement)
	}
	return strings.TrimSpace(body)
}
