// Package notify — TeeNotifier
//
// TeeNotifier mirrors every notification to exactly two downstream
// notifiers, guaranteeing both always receive the same message.
//
// # Behaviour
//
//   - Both notifiers are called on every Send, even if the first fails.
//   - If both return errors they are combined with errors.Join and
//     returned to the caller.
//   - If only one fails its error is returned; the other's success is
//     not affected.
//
// # Example
//
//	slack, _ := notify.NewSlackNotifier(slackURL)
//	log    := notify.NewLogNotifier(os.Stdout)
//	tee, _ := notify.NewTeeNotifier(slack, log)
//	// Every alert now goes to Slack AND the log simultaneously.
package notify
