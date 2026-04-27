// Package notify — QuotaNotifier
//
// QuotaNotifier enforces a hard cap on the total number of alert
// deliveries within a rolling time window. This is useful when the
// downstream channel (e.g. an SMS gateway or a paid API) imposes strict
// rate limits or cost constraints that go beyond simple per-path
// throttling.
//
// Usage:
//
//	inner, _ := notify.NewSlackNotifier(webhookURL)
//	q, err := notify.NewQuotaNotifier(inner, 100, time.Hour)
//	if err != nil { /* handle */ }
//
//	// Later:
//	remaining, resetsAt := q.Remaining()
//
// Once the quota is exhausted Send returns ErrQuotaExceeded. The window
// resets automatically; no manual intervention is required.
package notify
