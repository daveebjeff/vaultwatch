// Package notify – CheckpointNotifier
//
// CheckpointNotifier wraps any Notifier and maintains a per-path record of the
// last send attempt, its status, and whether it succeeded. This is useful for
// building status dashboards, health endpoints, or audit trails without adding
// external storage.
//
// Usage:
//
//	 inner := notify.NewLogNotifier(os.Stdout)
//	 cp, err := notify.NewCheckpointNotifier(inner)
//	 if err != nil { ... }
//
//	 // later, query the last known state of a secret path:
//	 rec, ok := cp.LastSeen("secret/db/password")
//	 if ok && !rec.Succeeded {
//	     log.Printf("last delivery to %s failed at %s", rec.Path, rec.SentAt)
//	 }
//
// Thread safety:
//
//	All methods are safe for concurrent use.
package notify
