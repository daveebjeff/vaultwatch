// Package notify — LabelNotifier
//
// LabelNotifier attaches static key/value metadata labels to every
// notification summary before forwarding it to an inner Notifier.
//
// Use it to tag alerts with environment, team, or service information
// so downstream systems (Slack channels, PagerDuty services, etc.) can
// route or filter them without modifying business logic.
//
// Example:
//
//	base := notify.NewLogNotifier(os.Stdout)
//	n, err := notify.NewLabelNotifier(base, map[string]string{
//		"env":  "production",
//		"team": "platform",
//	})
package notify
