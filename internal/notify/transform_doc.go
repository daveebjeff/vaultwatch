// # Transform Notifier
//
// The TransformNotifier applies a caller-supplied function to each
// [Message] before forwarding it to the inner [Notifier]. This is
// useful for enriching messages with additional metadata, redacting
// sensitive fields, or normalising body text in a pipeline.
//
// Example — upper-case all alert bodies:
//
//	base, _ := notify.NewLogNotifier(os.Stdout)
//	tn, err := notify.NewTransformNotifier(base, func(m notify.Message) notify.Message {
//		m.Body = strings.ToUpper(m.Body)
//		return m
//	})
//
// The transform function receives a copy of the message and must return
// the (possibly modified) message that will be forwarded.
package notify
