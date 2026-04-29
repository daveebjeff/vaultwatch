// Package notify – ExpBackoffNotifier
//
// ExpBackoffNotifier wraps any Notifier and retries failed deliveries
// using an exponential back-off strategy.
//
// # Behaviour
//
//   - The first attempt is made immediately.
//   - On failure the notifier waits initDelay before the next attempt.
//   - Each subsequent wait is multiplied by 2, capped at maxDelay.
//   - The context is respected between retries; cancellation stops
//     further attempts and returns ctx.Err().
//
// # Example
//
//	n, _ := notify.NewExpBackoffNotifier(
//		slack,
//		200*time.Millisecond,  // initDelay
//		30*time.Second,        // maxDelay
//		5,                     // max attempts
//	)
package notify
