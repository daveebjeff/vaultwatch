// Package notify – OnCallNotifier
//
// OnCallNotifier routes alert messages to the notifier that is currently
// on-call according to a set of time-bounded rotations.
//
// Example usage:
//
//	primary := notify.NewLogNotifier(os.Stdout)
//	backup  := notify.NewLogNotifier(os.Stderr)
//
//	rotations := []notify.OnCallRotation{
//		{
//			Name:     "week-team-a",
//			Start:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
//			End:      time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
//			Notifier: primary,
//		},
//		{
//			Name:     "week-team-b",
//			Start:    time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
//			End:      time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
//			Notifier: backup,
//		},
//	}
//
//	n, err := notify.NewOnCallNotifier(rotations)
//	if err != nil { /* handle */ }
//
// If the current time falls outside all rotation windows, Send returns
// ErrNoOnCallRotation so callers can decide how to handle gaps.
package notify
