// Package notify — BackpressureNotifier
//
// BackpressureNotifier decouples the caller from slow downstream notifiers by
// buffering outbound messages in a bounded channel. When the buffer is full the
// notifier rejects new messages immediately (returning ErrBackpressureQueueFull)
// rather than blocking the caller's goroutine.
//
// A background goroutine drains the queue sequentially, forwarding each message
// to the wrapped Notifier. On Stop(), the worker flushes any remaining queued
// messages before returning.
//
// Usage:
//
//	n, err := notify.NewBackpressureNotifier(inner, 256)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer n.Stop()
//
//	// Non-blocking — returns ErrBackpressureQueueFull when buffer is full.
//	if err := n.Send(ctx, msg); err != nil {
//	    log.Printf("dropped alert: %v", err)
//	}
package notify
