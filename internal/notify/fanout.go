package notify

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// FanoutNotifier sends a message to all inner notifiers concurrently,
// collecting any errors that occur.
type FanoutNotifier struct {
	notifiers []Notifier
}

// NewFanoutNotifier returns a FanoutNotifier that broadcasts to all provided
// notifiers in parallel. At least one notifier must be supplied.
func NewFanoutNotifier(notifiers ...Notifier) (*FanoutNotifier, error) {
	if len(notifiers) == 0 {
		return nil, fmt.Errorf("fanout: at least one notifier required")
	}
	for i, n := range notifiers {
		if n == nil {
			return nil, fmt.Errorf("fanout: notifier at index %d is nil", i)
		}
	}
	return &FanoutNotifier{notifiers: notifiers}, nil
}

// Send delivers msg to every inner notifier concurrently. All notifiers are
// attempted regardless of individual failures. A combined error is returned
// if any notifier fails.
func (f *FanoutNotifier) Send(ctx context.Context, msg Message) error {
	var (
		mu   sync.Mutex
		errs []string
		wg   sync.WaitGroup
	)

	for _, n := range f.notifiers {
		wg.Add(1)
		go func(n Notifier) {
			defer wg.Done()
			if err := n.Send(ctx, msg); err != nil {
				mu.Lock()
				errs = append(errs, err.Error())
				mu.Unlock()
			}
		}(n)
	}

	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("fanout: %d error(s): %s", len(errs), strings.Join(errs, "; "))
	}
	return nil
}
