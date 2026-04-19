package notify

import (
	"fmt"
	"net/http"
	"sync"
)

// PrometheusNotifier exposes alert counts as Prometheus metrics via an HTTP handler.
type PrometheusNotifier struct {
	mu       sync.Mutex
	counts   map[string]int
	inner    Notifier
}

// NewPrometheusNotifier wraps inner and tracks send counts by status.
func NewPrometheusNotifier(inner Notifier) (*PrometheusNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("prometheus: inner notifier must not be nil")
	}
	return &PrometheusNotifier{
		counts: make(map[string]int),
		inner:  inner,
	}, nil
}

// Send forwards the message to the inner notifier and increments the counter.
func (p *PrometheusNotifier) Send(msg Message) error {
	err := p.inner.Send(msg)
	p.mu.Lock()
	defer p.mu.Unlock()
	key := statusLabel(msg.Status)
	if err != nil {
		p.counts["send_error"]++
	} else {
		p.counts[key]++
	}
	return err
}

// Handler returns an http.HandlerFunc that serves basic Prometheus-style metrics.
func (p *PrometheusNotifier) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p.mu.Lock()
		defer p.mu.Unlock()
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		for label, count := range p.counts {
			fmt.Fprintf(w, "vaultwatch_notify_total{status=%q} %d\n", label, count)
		}
	}
}

func statusLabel(s Status) string {
	switch s {
	case StatusExpired:
		return "expired"
	case StatusExpiringSoon:
		return "expiring_soon"
	default:
		return "ok"
	}
}
