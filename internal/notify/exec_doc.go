// Package notify provides exec-based notifiers for vaultwatch.
//
// # HTTPGet Notifier
//
// HTTPGetNotifier fires an HTTP GET request to a configured base URL,
// appending alert details (status, secret path, expiry) as query parameters.
// Useful for simple webhook receivers or monitoring endpoints that accept GET.
//
// # Script Notifier
//
// ScriptNotifier executes a local script or binary when an alert fires,
// passing the status, secret path, and expiry timestamp as positional arguments.
// The script must be executable and return exit code 0 on success.
package notify
