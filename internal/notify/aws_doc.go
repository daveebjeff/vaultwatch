// Package notify provides AWS notifier implementations for VaultWatch.
//
// SNSNotifier publishes plain-text alert messages to an AWS SNS topic.
// It formats the message as "[STATUS] path — detail" and sets a subject
// derived from the alert status.
//
// SQSNotifier enqueues JSON-encoded notify.Message values to an AWS SQS
// queue, allowing downstream consumers to process alerts asynchronously.
//
// Both notifiers use the default AWS SDK credential chain (environment
// variables, shared credentials file, IAM role, etc.).
package notify
