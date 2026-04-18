// Package notify provides notifier implementations for VaultWatch alerts.
//
// Supported notifiers:
//   - LogNotifier: writes alerts to an io.Writer (default: stderr)
//   - SlackNotifier: posts alerts to a Slack incoming webhook
//   - EmailNotifier: sends alerts via SMTP email
//   - MultiNotifier: fans out alerts to multiple notifiers
//
// All notifiers implement the Notifier interface defined in notifier.go.
package notify
