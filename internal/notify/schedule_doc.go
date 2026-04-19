// Package notify — ScheduleNotifier
//
// ScheduleNotifier wraps any Notifier and suppresses messages that arrive
// outside of configured daily time windows. This is useful for silencing
// non-critical alerts overnight or during weekends.
//
// Example: only forward alerts between 08:00 and 18:00 UTC:
//
//	window := notify.TimeWindow{
//		Start: 8 * time.Hour,
//		End:   18 * time.Hour,
//	}
//	sn, err := notify.NewScheduleNotifier(inner, time.UTC, window)
//
// Multiple windows may be provided; a message is forwarded if it falls
// within ANY of the configured windows.
package notify
