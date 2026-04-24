package notify

import (
	"fmt"
	"time"
)

// EnrichNotifier wraps an inner Notifier and stamps each Message with
// additional computed labels before forwarding: a human-readable
// "time_to_expiry" string and a normalised "severity" label derived
// from the message Status.
type EnrichNotifier struct {
	inner Notifier
	now   func() time.Time // injectable for testing
}

// NewEnrichNotifier returns an EnrichNotifier that wraps inner.
// It returns an error if inner is nil.
func NewEnrichNotifier(inner Notifier) (*EnrichNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("enrich: inner notifier must not be nil")
	}
	return &EnrichNotifier{inner: inner, now: time.Now}, nil
}

// Send stamps msg.Labels with computed fields then delegates to the
// inner Notifier.  The original Message is not mutated; a shallow copy
// with a new Labels map is forwarded instead.
func (e *EnrichNotifier) Send(msg Message) error {
	enriched := msg
	enriched.Labels = make(map[string]string, len(msg.Labels)+2)
	for k, v := range msg.Labels {
		enriched.Labels[k] = v
	}

	if !msg.Expiry.IsZero() {
		ttl := msg.Expiry.Sub(e.now())
		if ttl < 0 {
			enriched.Labels["time_to_expiry"] = "expired"
		} else {
			enriched.Labels["time_to_expiry"] = ttl.Round(time.Second).String()
		}
	}

	enriched.Labels["severity"] = enrichSeverity(msg.Status)

	return e.inner.Send(enriched)
}

// enrichSeverity maps a Status value to a conventional severity string.
func enrichSeverity(s Status) string {
	switch s {
	case StatusExpired:
		return "critical"
	case StatusExpiringSoon:
		return "warning"
	default:
		return "info"
	}
}
