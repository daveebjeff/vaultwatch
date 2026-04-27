package notify

import (
	"fmt"
	"sync"
	"time"
)

// WatermarkNotifier forwards a message only when the expiry crosses a
// configured threshold boundary (e.g. first time TTL drops below 24h).
// Subsequent messages for the same path that remain below the threshold are
// suppressed until the secret is renewed above the watermark again.
type WatermarkNotifier struct {
	inner     Notifier
	watermark time.Duration

	mu      sync.Mutex
	below   map[string]bool // true = already fired for this threshold crossing
}

// NewWatermarkNotifier creates a WatermarkNotifier that fires once per
// threshold crossing. inner and a positive watermark duration are required.
func NewWatermarkNotifier(inner Notifier, watermark time.Duration) (*WatermarkNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("watermark: inner notifier must not be nil")
	}
	if watermark <= 0 {
		return nil, fmt.Errorf("watermark: duration must be positive, got %s", watermark)
	}
	return &WatermarkNotifier{
		inner:     inner,
		watermark: watermark,
		below:     make(map[string]bool),
	}, nil
}

// Send forwards msg to the inner notifier only on the first message where the
// remaining TTL crosses below the watermark. If the secret is later renewed
// (TTL rises above watermark), the state resets so the next crossing fires again.
func (w *WatermarkNotifier) Send(msg Message) error {
	timeLeft := time.Until(msg.Expiry)

	w.mu.Lock()
	alreadyFired := w.below[msg.Path]

	if timeLeft <= w.watermark {
		if alreadyFired {
			w.mu.Unlock()
			return nil // already notified for this crossing
		}
		w.below[msg.Path] = true
		w.mu.Unlock()
		return w.inner.Send(msg)
	}

	// TTL is above watermark — reset so next crossing fires again
	if alreadyFired {
		w.below[msg.Path] = false
	}
	w.mu.Unlock()
	return nil
}

// Reset clears the watermark state for all paths.
func (w *WatermarkNotifier) Reset() {
	w.mu.Lock()
	w.below = make(map[string]bool)
	w.mu.Unlock()
}
