package notify

import (
	"context"
	"fmt"
	"regexp"
)

// GrepNotifier forwards a message only when its body matches at least one
// of the provided regular expressions. It is the inverse of RedactNotifier:
// instead of scrubbing content it acts as a content-based gate.
type GrepNotifier struct {
	inner    Notifier
	patterns []*regexp.Regexp
}

// NewGrepNotifier returns a GrepNotifier that forwards messages whose body
// matches any of the supplied compiled patterns.
//
// Returns an error if inner is nil or no patterns are provided.
func NewGrepNotifier(inner Notifier, patterns []*regexp.Regexp) (*GrepNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("grep: inner notifier must not be nil")
	}
	if len(patterns) == 0 {
		return nil, fmt.Errorf("grep: at least one pattern is required")
	}
	for i, p := range patterns {
		if p == nil {
			return nil, fmt.Errorf("grep: pattern at index %d is nil", i)
		}
	}
	return &GrepNotifier{inner: inner, patterns: patterns}, nil
}

// Send forwards msg to the inner notifier only when msg.Body matches at
// least one pattern. Returns nil (without forwarding) when no pattern
// matches.
func (g *GrepNotifier) Send(ctx context.Context, msg Message) error {
	for _, p := range g.patterns {
		if p.MatchString(msg.Body) {
			return g.inner.Send(ctx, msg)
		}
	}
	return nil
}
