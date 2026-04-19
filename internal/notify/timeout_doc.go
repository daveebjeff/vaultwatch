// Package notify provides notification backends for vaultwatch.
//
// TimeoutNotifier
//
// TimeoutNotifier wraps any Notifier and enforces a hard deadline on each
// Send call. This is useful when integrating with slow or unreliable
// downstream systems where you do not want a single stalled notifier to
// block the monitoring loop.
//
// Example:
//
//	base, _ := notify.NewSlackNotifier(webhookURL)
//	tn, _ := notify.NewTimeoutNotifier(base, 5*time.Second)
//	// tn.Send will return an error if Slack does not respond within 5 s.
package notify
