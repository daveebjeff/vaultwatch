// Package notify provides notification backends for vaultwatch.
//
// FilterNotifier
//
// FilterNotifier wraps any Notifier and gates delivery on path prefixes.
// Only messages whose Path begins with one of the configured prefixes are
// forwarded to the inner notifier; all others are silently dropped.
//
// Usage:
//
//	 inner := notify.NewLogNotifier(os.Stdout)
//	 f, err := notify.NewFilterNotifier(inner, []string{"secret/prod/"})
//	 if err != nil { ... }
//	 f.Send(msg)
package notify
