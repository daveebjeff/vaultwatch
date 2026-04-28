// Package notify provides the CoalesceNotifier, which merges rapid-fire
// notifications for the same secret path into a single delivery.
//
// # Overview
//
// When a secret is nearing expiry, the monitor may emit several alerts in quick
// succession — for example, one per check interval during a busy polling window.
// CoalesceNotifier absorbs those duplicates and forwards only the most-recent
// message once a configurable quiet window has elapsed with no new arrivals for
// the same path.
//
// This is similar to DebounceNotifier but is optimised for the case where the
// caller wants the *latest* state of a path rather than the first.
//
// # Behaviour
//
//   - Each unique secret path has its own independent timer.
//   - When a message arrives the timer is reset; the message is stored.
//   - When the timer fires (no new message for the window) the stored message
//     is forwarded to the inner Notifier.
//   - A manual Flush call forwards all pending messages immediately and resets
//     all timers.
//
// # Configuration
//
//	notifier, err := notify.NewCoalesceNotifier(inner,
//	    notify.WithCoalesceWindow(10 * time.Second),
//	    notify.WithCoalesceMaxPending(100),
//	)
//
// # When to use
//
// Use CoalesceNotifier when:
//   - Your polling interval is short and you want to avoid alert storms.
//   - You only care about the latest status of a path, not every intermediate
//     state.
//   - You are upstream of a rate-limited channel such as PagerDuty or SMS.
package notify
