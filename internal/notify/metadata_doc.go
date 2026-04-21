// Package notify — MetadataNotifier
//
// MetadataNotifier enriches every outgoing alert Message with a static set of
// key-value labels before forwarding it to the wrapped Notifier.  This is
// useful for stamping environment, team, or datacenter context onto all
// notifications produced by a particular pipeline without modifying the
// upstream monitor configuration.
//
// # Behaviour
//
//   - Metadata keys are merged into Message.Labels before the message is
//     forwarded.  Keys that are already present in the message are left
//     unchanged, so per-message labels always take precedence.
//   - The metadata map is safe for concurrent reads and can be replaced at
//     runtime via SetMetadata.
//
// # Example
//
//	base := notify.NewSlackNotifier(webhookURL)
//	n, err := notify.NewMetadataNotifier(base, map[string]string{
//	    "env":    "production",
//	    "region": "us-east-1",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
package notify
