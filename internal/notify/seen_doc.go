// Package notify provides the SeenNotifier, which suppresses repeated
// notifications for a given secret path within a configurable rolling time
// window.
//
// # Behaviour
//
// The first message for a path is always forwarded to the inner notifier and
// the arrival time is recorded. Any subsequent message for the same path is
// silently dropped until the configured window has fully elapsed since that
// first arrival. Once the window expires, the next message is forwarded and
// the clock is reset.
//
// This is distinct from DedupNotifier (which keys on status changes) and
// CachedNotifier (which keys on message content). SeenNotifier is purely
// time-based and path-scoped.
//
// # Usage
//
//	sn, err := notify.NewSeenNotifier(inner, 4*time.Hour)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// # Manual control
//
// Call Forget(path) to immediately allow the next message for a specific path
// to be forwarded, or Reset() to clear all tracked paths at once.
package notify
