// Package notify provides notification backends for vaultwatch.
//
// # RollupNotifier
//
// RollupNotifier batches multiple secret expiration alerts into a single
// summary message, reducing notification noise when many secrets expire
// around the same time.
//
// Messages are held in memory until either:
//   - The configured time window elapses (timer-based flush), or
//   - The number of buffered messages reaches maxSize.
//
// The summary message carries the worst Status observed across all
// buffered messages, making it easy to triage in downstream systems.
//
// Example usage:
//
//	base := notify.NewLogNotifier(os.Stdout)
//	r, err := notify.NewRollupNotifier(base, 5*time.Minute, 20)
package notify
