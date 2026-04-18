// Package notify provides notifier implementations for alerting on
// Vault secret expiration and lease renewal events.
//
// HTTP Notifiers
//
// HTTPGetNotifier sends a simple GET request to a configured URL.
//
// WebhookNotifier sends a JSON POST payload to a configured endpoint.
//
// Both notifiers are useful for triggering lightweight integrations
// or custom automation pipelines without requiring a full webhook
// receiver setup.
package notify
