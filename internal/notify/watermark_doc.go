// Package notify — WatermarkNotifier
//
// WatermarkNotifier fires exactly once each time a secret's remaining TTL
// crosses below a configured threshold (the "watermark").
//
// Behaviour
//
//   - When the TTL first drops at or below the watermark, the message is
//     forwarded to the inner notifier and the path is marked as "fired".
//   - Further messages for the same path that remain below the watermark are
//     suppressed — preventing alert storms during repeated monitor ticks.
//   - If the secret is subsequently renewed and the TTL rises above the
//     watermark, the fired state is cleared. The next threshold crossing will
//     fire again.
//
// Example
//
//	win, _ := notify.NewWatermarkNotifier(slackNotifier, 24*time.Hour)
//	// Fires once when a secret has less than 24 h remaining, then again
//	// only after it has been renewed past 24 h and drops below once more.
package notify
