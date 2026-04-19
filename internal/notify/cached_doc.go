// Package notify — CachedNotifier
//
// CachedNotifier wraps any Notifier and suppresses repeated alerts for
// the same secret path and status within a configurable time-to-live
// (TTL) window.
//
// Use it to prevent alert storms when a secret is expiring and the
// monitor loop fires repeatedly before the secret is renewed.
//
// Example:
//
//	base, _ := NewSlackNotifier(webhookURL)
//	cached, _ := NewCachedNotifier(base, 30*time.Minute)
//	// Only the first alert per path+status is forwarded within 30 min.
package notify
