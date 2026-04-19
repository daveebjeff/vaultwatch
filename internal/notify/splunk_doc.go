// Package notify provides notification backends for vaultwatch alerts.
//
// # Splunk Notifier
//
// SplunkNotifier delivers alert messages to a Splunk HTTP Event Collector
// (HEC) endpoint. Each message is wrapped in a standard HEC JSON envelope
// with the sourcetype set to "vaultwatch".
//
// Usage:
//
//	n, err := notify.NewSplunkNotifier(
//		"https://splunk.example.com:8088/services/collector/event",
//		"<hec-token>",
//	)
//
The HEC token is sent via the Authorization header as "Splunk <token>".
package notify
