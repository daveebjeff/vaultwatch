// Package notify provides notification primitives for vaultwatch.
//
// SnapshotNotifier
//
// SnapshotNotifier wraps any Notifier and maintains an in-memory record of
// the most recently delivered message for each secret path.
//
// Use Latest to retrieve the current state of a single path, or All to
// obtain a full map of path → Snapshot for status-page or health-endpoint
// use cases.
//
// Example:
//
//	base := notify.NewLogNotifier(os.Stdout)
//	sn, _ := notify.NewSnapshotNotifier(base)
//
//	// later…
//	if snap, ok := sn.Latest("secret/db/password"); ok {
//		fmt.Println(snap.Message.Status, snap.ReceivedAt)
//	}
package notify
