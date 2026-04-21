// Package notify — ReplayNotifier
//
// ReplayNotifier wraps any Notifier and maintains an in-memory ring of recent
// messages. It is designed for scenarios where a newly registered downstream
// target needs to receive recent alert history without requiring a full
// re-scan of Vault.
//
// Usage:
//
//	r, err := notify.NewReplayNotifier(inner, 100, 24*time.Hour)
//	if err != nil { ... }
//
//	// Normal operation — messages flow through to inner.
//	r.Send(msg)
//
//	// Later, replay retained messages to a new subscriber.
//	r.Replay(newSubscriber)
//
// Retention policy:
//   - Up to maxItems messages are kept (oldest evicted first).
//   - Messages whose ExpiresAt timestamp is older than maxAge are pruned
//     automatically on each Send or Len call.
//
// Thread safety: all methods are safe for concurrent use.
package notify
