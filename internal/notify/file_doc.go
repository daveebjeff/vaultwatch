// Package notify provides notifier implementations for vaultwatch alerts.
//
// FileNotifier writes one JSON-encoded Message per line to a local file,
// suitable for log aggregation pipelines (e.g. Fluentd, Filebeat).
//
// SyslogNotifier forwards alerts to the local syslog daemon using the
// standard log/syslog package. Severity is chosen by alert status:
//
//	StatusExpired      → syslog.Err
//	StatusExpiringSoon → syslog.Warning
//	other              → syslog.Info
//
// Both notifiers implement the Notifier interface and can be composed with
// MultiNotifier to fan out to multiple destinations simultaneously.
//
// # Notifier Interface
//
// All notifiers in this package satisfy the following interface:
//
//	type Notifier interface {
//		Notify(ctx context.Context, msg Message) error
//		Close() error
//	}
//
// Implementations should be safe for concurrent use unless otherwise noted.
package notify
