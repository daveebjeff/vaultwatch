package notify

import (
	"context"
	"sync"
)

// TapNotifier forwards every message to the inner notifier and also calls
// a user-supplied tap function with a copy of the message. The tap function
// runs synchronously before the inner Send returns. Errors from the tap
// function are silently discarded so that they never affect the primary
// notification path.
//
// A typical use-case is capturing messages in tests or feeding a secondary
// pipeline without changing the observable behaviour of the wrapped notifier.
type TapNotifier struct {
	mu    sync.Mutex
	inner Notifier
	tapFn func(Message)
}

// NewTapNotifier wraps inner and calls tapFn for every message that passes
// through. Both inner and tapFn must be non-nil.
func NewTapNotifier(inner Notifier, tapFn func(Message)) (*TapNotifier, error) {
	if inner == nil {
		return nil, ErrNilInner
	}
	if tapFn == nil {
		return nil, errNilTapFn
	}
	return &TapNotifier{inner: inner, tapFn: tapFn}, nil
}

// Send calls the tap function with a copy of msg and then forwards msg to
// the inner notifier. The error from the inner notifier is returned; any
// panic inside tapFn is recovered and suppressed.
func (t *TapNotifier) Send(ctx context.Context, msg Message) error {
	t.mu.Lock()
	fn := t.tapFn
	t.mu.Unlock()

	// Run the tap function, suppressing any panic so the primary path is safe.
	func() {
		defer func() { recover() }() //nolint:errcheck
		fn(msg)
	}()

	return t.inner.Send(ctx, msg)
}

// SetTapFn replaces the tap function at runtime. It is safe to call
// concurrently with Send.
func (t *TapNotifier) SetTapFn(fn func(Message)) {
	if fn == nil {
		return
	}
	t.mu.Lock()
	t.tapFn = fn
	t.mu.Unlock()
}
