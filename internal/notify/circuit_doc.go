// Package notify provides circuit breaker support via CircuitNotifier.
//
// CircuitNotifier wraps any Notifier and prevents cascading failures by
// tracking consecutive errors. Once the failure count reaches maxFailures
// the circuit opens and all subsequent Send calls are rejected immediately
// without calling the inner notifier.
//
// After the resetAfter duration the circuit transitions to half-open and
// allows a single probe. A successful probe closes the circuit; a failed
// probe re-opens it.
//
// Example:
//
//	base := notify.NewSlackNotifier(webhookURL)
//	cb, err := notify.NewCircuitNotifier(base, 5, 30*time.Second)
//	if err != nil {
//		log.Fatal(err)
//	}
package notify
