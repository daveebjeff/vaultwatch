// Package notify — DebounceNotifier
//
// DebounceNotifier suppresses rapid-fire notifications for the same secret
// path by holding back delivery until a quiet period has elapsed. This is
// useful when a secret is close to expiry and the monitor loop fires many
// alerts in quick succession before the secret is renewed.
//
// Usage:
//
//	base, _ := notify.NewLogNotifier(os.Stdout)
//	d, err := notify.NewDebounceNotifier(base, 10*time.Second)
//	if err != nil { ... }
//
//	// Rapid sends for the same path — only the last one is forwarded
//	// after 10 s of silence.
//	d.Send(msg1)
//	d.Send(msg2)
//	d.Send(msg3) // only msg3 reaches base
//
// Notes:
//   - Each unique Message.Path has its own independent timer.
//   - The notifier is safe for concurrent use.
//   - Errors from the inner notifier are silently discarded because delivery
//     happens asynchronously in a timer goroutine.
package notify
