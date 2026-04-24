// Package notify — TraceIDNotifier
//
// TraceIDNotifier stamps every outgoing Message with a unique trace ID so
// that a single alert event can be correlated across multiple notifiers,
// log lines, and external systems.
//
// Usage:
//
//	base := notify.NewLogNotifier(os.Stdout)
//	n, err := notify.NewTraceIDNotifier(base, "")
//	if err != nil {
//		log.Fatal(err)
//	}
//	// Each Send call will have msg.Labels["trace_id"] set automatically.
//
// Propagating an existing trace ID:
//
//	ctx := notify.ContextWithTraceID(context.Background(), "abc123")
//	_ = n.Send(ctx, msg) // msg.Labels["trace_id"] == "abc123"
//
// The label key can be customised:
//
//	n, _ := notify.NewTraceIDNotifier(base, "x-request-id")
package notify
