package notify

import (
	"context"
	"fmt"
	"sync"
)

// BackpressureNotifier wraps a Notifier and applies a bounded in-memory queue.
// When the queue is full, Send returns ErrBackpressureQueueFull instead of
// blocking the caller. A background worker drains the queue sequentially.
type BackpressureNotifier struct {
	inner    Notifier
	queue    chan Message
	stopOnce sync.Once
	stop     chan struct{}
	wg       sync.WaitGroup
}

// NewBackpressureNotifier creates a BackpressureNotifier with the given queue
// capacity. The background drain goroutine starts immediately.
func NewBackpressureNotifier(inner Notifier, capacity int) (*BackpressureNotifier, error) {
	if inner == nil {
		return nil, ErrBackpressureNilInner
	}
	if capacity <= 0 {
		return nil, ErrBackpressureZeroCapacity
	}
	n := &BackpressureNotifier{
		inner: inner,
		queue: make(chan Message, capacity),
		stop:  make(chan struct{}),
	}
	n.wg.Add(1)
	go n.drain()
	return n, nil
}

// Send enqueues msg for async delivery. Returns ErrBackpressureQueueFull when
// the internal buffer is at capacity.
func (n *BackpressureNotifier) Send(ctx context.Context, msg Message) error {
	select {
	case n.queue <- msg:
		return nil
	default:
		return fmt.Errorf("%w: capacity %d", ErrBackpressureQueueFull, cap(n.queue))
	}
}

// Stop signals the background worker to exit and waits for it to finish.
func (n *BackpressureNotifier) Stop() {
	n.stopOnce.Do(func() { close(n.stop) })
	n.wg.Wait()
}

func (n *BackpressureNotifier) drain() {
	defer n.wg.Done()
	for {
		select {
		case msg := <-n.queue:
			//nolint:errcheck — delivery errors are best-effort in async mode
			_ = n.inner.Send(context.Background(), msg)
		case <-n.stop:
			// Drain remaining items before exiting.
			for {
				select {
				case msg := <-n.queue:
					_ = n.inner.Send(context.Background(), msg)
				default:
					return
				}
			}
		}
	}
}
