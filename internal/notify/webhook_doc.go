// Package notify provides notification delivery for VaultWatch alerts.
//
// # Webhook Notifier
//
// The WebhookNotifier sends alert messages as HTTP POST requests with a
// JSON body to a configurable URL endpoint.
//
// Example:
//
//	n, err := notify.NewWebhookNotifier("https://example.com/hook")
//	if err != nil {
//		log.Fatal(err)
//	}
//
// The payload includes the secret path, status, expiry time, and message.
package notify
