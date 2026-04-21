package notify

import (
	"sync"
	"time"
)

// RequeueNotifier holds failed messages and retries them on the next Send call
// or when Flush is called explicitly. It is useful when you want to preserve
// failed notifications and attempt redelivery without losing the original event.
type RequeueNotifier struct {
	inner    Notifier
	mu       sync.Mutex
	queue    []Message
	maxQueue int
	retryAge time.Duration
}

// NewRequeueNotifier wraps inner and buffers up to maxQueue failed messages.
// retryAge is the minimum time a message must wait before being retried.
func NewRequeueNotifier(inner Notifier, maxQueue int, retryAge time.Duration) (*RequeueNotifier, error) {
	if inner == nil {
		return nil, ErrNilInner
	}
	if maxQueue <= 0 {
		return nil, ErrZeroMax
	}
	if retryAge <= 0 {
		return nil, ErrZeroDuration
	}
	return &RequeueNotifier{
		inner:    inner,
		maxQueue: maxQueue,
		retryAge:  retryAge,
	}, nil
}

// Send first attempts to flush any queued messages that have aged out, then
// delivers msg. If delivery of msg fails it is appended to the queue (up to
// maxQueue). Oldest entries are dropped when the queue is full.
func (r *RequeueNotifier) Send(msg Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.flushQueued()

	if err := r.inner.Send(msg); err != nil {
		r.enqueue(msg)
		return err
	}
	return nil
}

// Flush attempts to deliver all queued messages regardless of age.
func (r *RequeueNotifier) Flush() {
	r.mu.Lock()
	defer r.mu.Unlock()

	remaining := r.queue[:0]
	for _, m := range r.queue {
		if err := r.inner.Send(m); err != nil {
			remaining = append(remaining, m)
		}
	}
	r.queue = remaining
}

// QueueLen returns the number of messages currently buffered.
func (r *RequeueNotifier) QueueLen() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.queue)
}

// flushQueued retries messages whose queued timestamp has exceeded retryAge.
// Must be called with r.mu held.
func (r *RequeueNotifier) flushQueued() {
	cutoff := time.Now().Add(-r.retryAge)
	remaining := r.queue[:0]
	for _, m := range r.queue {
		if m.At.Before(cutoff) {
			if err := r.inner.Send(m); err != nil {
				remaining = append(remaining, m)
			}
		} else {
			remaining = append(remaining, m)
		}
	}
	r.queue = remaining
}

// enqueue appends msg, dropping the oldest entry if the queue is full.
func (r *RequeueNotifier) enqueue(msg Message) {
	if len(r.queue) >= r.maxQueue {
		r.queue = r.queue[1:]
	}
	r.queue = append(r.queue, msg)
}
