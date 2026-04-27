// Package notify – CooldownNotifier
//
// CooldownNotifier wraps any Notifier and enforces a per-path quiet period
// after each successfully forwarded message. While the cooldown is active,
// subsequent sends for the same path are silently dropped.
//
// This differs from RateLimitNotifier (fixed window) and ThrottleNotifier
// (max count per window) by providing a simple "don't fire again for N seconds
// after the last fire" semantic, which is useful for noisy secrets that flip
// between expiring and healthy states rapidly.
//
// Usage:
//
//	n, _ := notify.NewCooldownNotifier(inner, 5*time.Minute)
//	_ = n.Send(ctx, msg)   // forwarded
//	_ = n.Send(ctx, msg)   // suppressed – cooldown active
//
package notify
