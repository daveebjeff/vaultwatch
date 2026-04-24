package notify

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// traceIDKey is the context key used to store and retrieve a trace ID.
type traceIDKey struct{}

// TraceIDNotifier wraps an inner Notifier and stamps each outgoing Message
// with a unique trace ID stored in Message.Labels["trace_id"]. If the
// context already carries a trace ID it is reused, enabling correlation
// across a chain of notifiers.
type TraceIDNotifier struct {
	inner  Notifier
	header string // label key, defaults to "trace_id"
}

// NewTraceIDNotifier returns a TraceIDNotifier that decorates inner.
// header is the label key written to each message; pass an empty string to
// use the default "trace_id".
func NewTraceIDNotifier(inner Notifier, header string) (*TraceIDNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("traceid: inner notifier must not be nil")
	}
	if header == "" {
		header = "trace_id"
	}
	return &TraceIDNotifier{inner: inner, header: header}, nil
}

// Send stamps msg with a trace ID then forwards it to the inner notifier.
func (t *TraceIDNotifier) Send(ctx context.Context, msg Message) error {
	id := traceIDFromContext(ctx)
	if id == "" {
		id = newTraceID()
	}
	if msg.Labels == nil {
		msg.Labels = make(map[string]string)
	}
	msg.Labels[t.header] = id
	return t.inner.Send(ctx, msg)
}

// ContextWithTraceID returns a new context carrying the given trace ID.
func ContextWithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, id)
}

func traceIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(traceIDKey{}).(string)
	return v
}

func newTraceID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
