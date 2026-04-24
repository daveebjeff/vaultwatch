// Package notify — EnrichNotifier
//
// EnrichNotifier is a middleware notifier that automatically stamps
// outgoing messages with computed observability labels before they
// reach the downstream notifier.
//
// Labels added
//
//   - time_to_expiry  Human-readable duration until the secret expires
//                     (e.g. "4h30m0s"), or the string "expired" when the
//                     expiry time is already in the past.  The label is
//                     omitted when Message.Expiry is the zero value.
//
//   - severity        Conventional severity string derived from the
//                     message Status:
//                       StatusExpired      → "critical"
//                       StatusExpiringSoon → "warning"
//                       otherwise          → "info"
//
// Usage
//
//	base, _ := notify.NewSlackNotifier(webhookURL)
//	en, _ := notify.NewEnrichNotifier(base)
//	// en now forwards messages with extra labels to Slack.
package notify
