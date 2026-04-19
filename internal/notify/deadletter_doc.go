// Package notify provides notification primitives for vaultwatch.
//
// # DeadLetterNotifier
//
// DeadLetterNotifier wraps any Notifier and captures messages whose delivery
// failed so they can be inspected or retried by the caller.
//
// Usage:
//
//	base := notify.NewSlackNotifier(webhookURL)
//	dl, err := notify.NewDeadLetterNotifier(base, 200)
//	if err != nil { ... }
//
//	// later, inspect failures
//	for _, entry := range dl.Drain() {
//	    log.Printf("failed %s at %s: %v", entry.Message.Path, entry.FailedAt, entry.Err)
//	}
//
// The queue is bounded by maxSize; once full, additional failures are silently
// dropped to avoid unbounded memory growth. Use Drain to consume and clear the
// queue, or Failed to read without clearing.
package notify
