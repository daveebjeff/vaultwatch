package notify

import (
	"context"
	"sync/atomic"
)

// mockNotifier is a test double that delegates Send to an arbitrary function.
type mockNotifier struct {
	fn      func(ctx context.Context, msg Message) error
	calls   atomic.Int64
}

func (m *mockNotifier) Send(ctx context.Context, msg Message) error {
	m.calls.Add(1)
	if m.fn != nil {
		return m.fn(ctx, msg)
	}
	return nil
}
