// Package notify provides notification backends for VaultWatch alerts.
//
// All backends implement the Notifier interface and accept a Message value
// that describes the secret path, its expiry status and a human-readable
// detail string.
//
// Available backends:
//   - LogNotifier   – writes structured log lines to an io.Writer
//   - SlackNotifier – posts messages to a Slack incoming webhook
//   - EmailNotifier – sends SMTP email alerts
//   - PagerDutyNotifier – triggers PagerDuty incidents via Events API v2
//   - WebhookNotifier   – posts JSON payloads to an arbitrary HTTP endpoint
//   - OpsGenieNotifier  – creates OpsGenie alerts via the REST API
//   - MultiNotifier    – fans out to multiple Notifier implementations
package notify
