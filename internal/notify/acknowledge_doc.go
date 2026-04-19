// Package notify provides the AcknowledgeNotifier, which wraps any Notifier
// and suppresses repeated alert delivery for secret paths that an operator
// has explicitly acknowledged.
//
// Once a path is acknowledged, alerts for that path are silently dropped
// until the configured TTL elapses. This is useful in on-call workflows
// where an engineer has seen an alert and wants to prevent further noise
// while they work on the issue.
//
// Example usage:
//
//	base := notify.NewSlackNotifier(webhookURL)
//	an, _ := notify.NewAcknowledgeNotifier(base, 4*time.Hour)
//
//	// Operator acknowledges a noisy secret path for 4 hours:
//	an.Acknowledge("secret/prod/db")
//
//	// Subsequent Send calls for that path are suppressed until TTL expires.
package notify
