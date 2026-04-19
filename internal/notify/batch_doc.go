// Package notify — BatchNotifier
//
// BatchNotifier accumulates alert messages over a configurable time window
// or until a maximum batch size is reached, then forwards a single summary
// message to the wrapped Notifier.
//
// This is useful for reducing noise when many secrets expire simultaneously.
//
// Usage:
//
//	batch, err := notify.NewBatchNotifier(inner, 30*time.Second, 20)
//	if err != nil { ... }
//	// Call batch.Flush(ctx) on shutdown to drain pending messages.
package notify
