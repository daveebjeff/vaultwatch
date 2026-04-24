// Package notify – WindowNotifier
//
// WindowNotifier enforces a true sliding-window rate limit across all messages
// regardless of their path or status.  It complements ThrottleNotifier, which
// uses fixed reset periods, by providing smoother, burst-resistant control.
//
// Usage:
//
//	base := notify.NewLogNotifier(os.Stdout)
//	wn, err := notify.NewWindowNotifier(base, 5, 30*time.Second)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// At most 5 messages are forwarded in any rolling 30-second window.
//
// When the limit is exceeded Send returns notify.ErrSuppressed and the inner
// notifier is not called.
package notify
