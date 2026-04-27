// Package notify — HedgeNotifier
//
// HedgeNotifier implements a hedged-request pattern for notifications. When a
// primary notifier is slow, a secondary is launched in parallel after a
// configurable delay so that the overall latency is bounded by whichever
// backend responds first.
//
// Usage:
//
//	primary := notify.NewSlackNotifier(slackURL)
//	backup  := notify.NewWebhookNotifier(backupURL)
//	h, err  := notify.NewHedgeNotifier(primary, backup, 500*time.Millisecond)
//
// The hedge delay should be chosen to be slightly above the p99 latency of the
// primary notifier so that the secondary is only invoked when the primary is
// genuinely slow or failing.
//
// Both notifiers receive identical Message values. If either returns nil the
// call succeeds immediately. If both fail, the primary error is surfaced.
package notify
