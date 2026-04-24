// Package notify — InspectNotifier
//
// InspectNotifier is a transparent middleware that records every notification
// event passing through it. It is useful in two scenarios:
//
//  1. Testing — wrap a NoopNotifier (or any notifier) with InspectNotifier to
//     assert which messages were sent and what errors, if any, were returned.
//
//  2. Debugging — insert InspectNotifier into a live pipeline to log or inspect
//     messages without altering the delivery behaviour.
//
// Example:
//
//	inner, _ := notify.NewSlackNotifier(webhookURL)
//	inspector, _ := notify.NewInspectNotifier(inner)
//
//	// use inspector in your pipeline …
//
//	// later, inspect what was sent:
//	for _, e := range inspector.Entries() {
//		fmt.Println(e.Message.Path, e.Err)
//	}
package notify
