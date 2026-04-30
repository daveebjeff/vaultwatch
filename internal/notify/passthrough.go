package notify

import (
	"context"
	"fmt"
)

// PassthroughNotifier forwards every message to the inner notifier unchanged,
// but records the total number of messages seen and sent. It is useful as a
// lightweight instrumentation shim when you want raw send counts without the
// overhead of the full InspectNotifier.
type PassthroughNotifier struct {
	inner   Notifier
	seen    int64
	sent    int64
	errors  int64
}

// NewPassthroughNotifier wraps inner with pass-through counting.
// It returns an error if inner is nil.
func NewPassthroughNotifier(inner Notifier) (*PassthroughNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("passthrough: inner notifier must not be nil")
	}
	return &PassthroughNotifier{inner: inner}, nil
}

// Send forwards msg to the inner notifier and updates counters.
func (p *PassthroughNotifier) Send(ctx context.Context, msg Message) error {
	p.seen++
	err := p.inner.Send(ctx, msg)
	if err != nil {
		p.errors++
		return err
	}
	p.sent++
	return nil
}

// Seen returns the total number of messages passed to Send.
func (p *PassthroughNotifier) Seen() int64 { return p.seen }

// Sent returns the number of messages successfully forwarded.
func (p *PassthroughNotifier) Sent() int64 { return p.sent }

// Errors returns the number of messages that resulted in an error.
func (p *PassthroughNotifier) Errors() int64 { return p.errors }

// Reset zeroes all counters.
func (p *PassthroughNotifier) Reset() {
	p.seen = 0
	p.sent = 0
	p.errors = 0
}
