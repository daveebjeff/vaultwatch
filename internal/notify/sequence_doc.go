// Package notify — SequenceNotifier
//
// SequenceNotifier delivers a notification through an ordered chain of
// notifiers, stopping immediately on the first failure. This is useful when
// steps have dependencies — for example, enriching a message before writing
// it to a file and then forwarding it to an external webhook.
//
// Contrast with MultiNotifier, which attempts every notifier regardless of
// intermediate errors, and FanoutNotifier, which runs all notifiers
// concurrently. SequenceNotifier is strictly serial and fail-fast.
//
// Example:
//
//	sn, err := notify.NewSequenceNotifier(
//		enricher,   // step 0: attach severity labels
//		fileWriter, // step 1: persist to disk
//		webhook,    // step 2: forward to external system
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
package notify
