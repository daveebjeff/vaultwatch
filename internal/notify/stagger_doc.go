// Package notify — StaggerNotifier
//
// StaggerNotifier delivers a notification to a series of inner notifiers
// with a configurable delay inserted between each successive send. This is
// useful when you have multiple downstream services that should all receive
// an alert but cannot absorb a burst of simultaneous requests.
//
// Example usage:
//
//	// Notify three channels with 250 ms between each delivery.
//	s, err := notify.NewStaggerNotifier(
//		250*time.Millisecond,
//		slackNotifier,
//		pagerDutyNotifier,
//		emailNotifier,
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Behaviour:
//   - Notifiers are called in the order they were provided.
//   - The delay is applied after every notifier except the last.
//   - If the context is cancelled during a delay the send is aborted
//     immediately and the context error is returned.
//   - The first notifier error encountered is returned; remaining notifiers
//     are still attempted so no channel is silently skipped.
//   - New notifiers can be added at runtime via Add.
package notify
