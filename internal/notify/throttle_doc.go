// Package notify provides ThrottleNotifier, which limits the total number of
// notifications forwarded to an inner Notifier within a sliding time window.
//
// Once the configured maximum count is reached for the current window, all
// further Send calls are silently dropped until the window expires and resets.
//
// Example usage:
//
//	base := notify.NewSlackNotifier(webhookURL)
//	th, err := notify.NewThrottleNotifier(base, 10, time.Hour)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// th will forward at most 10 messages per hour to Slack.
package notify
