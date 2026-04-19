// Package notify includes notifiers for observability and monitoring platforms.
//
// Available monitoring/SIEM notifiers:
//
//   - DatadogNotifier  — posts events to the Datadog Events API
//   - SplunkNotifier   — sends events to a Splunk HEC endpoint
//
// These notifiers are useful for centralising vault secret expiry signals
// into existing observability pipelines alongside infrastructure metrics and
// logs.
//
// All monitoring notifiers accept a notify.Message and translate the Status,
// Path, ExpiresAt, and Body fields into the platform-specific payload format.
package notify
