// Package notify provides notification delivery for vaultwatch.
//
// # Escalation
//
// EscalationNotifier wraps a primary and secondary [Notifier].
// Alerts are first sent to the primary. If the primary returns an error the
// secondary is tried immediately as a fallback.
//
// When the primary succeeds the alert is held in a pending set. Callers should
// invoke [EscalationNotifier.Escalate] periodically (e.g. on each monitor
// tick). Any alert that has remained unacknowledged for longer than the
// configured timeout is forwarded to the secondary and removed from the
// pending set.
//
// Use [EscalationNotifier.Acknowledge] to mark an alert resolved so it is not
// escalated even if the timeout has not yet elapsed.
package notify
