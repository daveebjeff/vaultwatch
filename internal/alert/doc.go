// Package alert provides evaluation logic and notification interfaces
// for Vault secret expiration alerts. It determines alert severity based
// on a configurable warning threshold and supports pluggable notifiers.
//
// Severity levels:
//   - Critical: the secret has already expired or expires within half the warning threshold.
//   - Warning:  the secret expires within the configured warning threshold.
//   - OK:       the secret is not approaching expiration.
//
// Notifiers implement the Notifier interface and can be composed to fan
// out alerts to multiple backends (e.g. Slack, PagerDuty, email).
package alert
