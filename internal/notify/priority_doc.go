// Package notify — PriorityNotifier
//
// PriorityNotifier routes each notification to the most appropriate registered
// notifier based on the severity of the message status.
//
// # Priority levels
//
// Each notifier is registered at an integer priority level. The notifier with
// the highest level that is greater than or equal to the message's severity
// is selected. Severity is derived from Status:
//
//	  StatusExpired      → 2
//	  StatusExpiringSoon → 1
//	  StatusOK           → 0
//
// # Example
//
//	 p := notify.NewPriorityNotifier()
//	 _ = p.Add(0, logNotifier)       // catch-all for OK
//	 _ = p.Add(1, slackNotifier)     // warn on expiring-soon
//	 _ = p.Add(2, pagerdutyNotifier) // page on expired
//
// When a StatusExpired message arrives the PagerDuty notifier is used.
// A StatusExpiringSoon message is routed to Slack.
// A StatusOK message falls through to the log notifier.
package notify
