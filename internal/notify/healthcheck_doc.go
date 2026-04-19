// Package notify provides notification primitives for vaultwatch.
//
// # HealthCheck Notifier
//
// NewHealthCheckNotifier wraps any Notifier and tracks its operational health.
// A background goroutine periodically sends a synthetic probe message so that
// health degrades automatically when the downstream target becomes unreachable,
// without requiring a real secret event to trigger the check.
//
// Usage:
//
//	h, err := notify.NewHealthCheckNotifier(slackNotifier, "slack", 30*time.Second)
//	if err != nil { ... }
//	defer h.Stop()
//
//	// later:
//	if !h.Status().Healthy {
//		log.Printf("slack notifier unhealthy: %v", h.Status().LastError)
//	}
package notify
