package notify

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// HealthStatus represents the health of a notifier.
type HealthStatus struct {
	Name      string
	Healthy   bool
	LastError error
	CheckedAt time.Time
}

// HealthCheckNotifier wraps a Notifier and periodically probes it with a
// synthetic message, exposing the result via Status().
type HealthCheckNotifier struct {
	inner    Notifier
	name     string
	interval time.Duration
	mu       sync.RWMutex
	status   HealthStatus
	stop     chan struct{}
}

// NewHealthCheckNotifier creates a HealthCheckNotifier.
// interval is how often the probe runs in the background.
func NewHealthCheckNotifier(inner Notifier, name string, interval time.Duration) (*HealthCheckNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("healthcheck: inner notifier must not be nil")
	}
	if interval <= 0 {
		return nil, fmt.Errorf("healthcheck: interval must be positive")
	}
	h := &HealthCheckNotifier{
		inner:    inner,
		name:     name,
		interval: interval,
		stop:     make(chan struct{}),
		status:   HealthStatus{Name: name, Healthy: true, CheckedAt: time.Now()},
	}
	go h.loop()
	return h, nil
}

// Send forwards the message to the inner notifier.
func (h *HealthCheckNotifier) Send(ctx context.Context, msg Message) error {
	err := h.inner.Send(ctx, msg)
	h.record(err)
	return err
}

// Status returns the latest health status.
func (h *HealthCheckNotifier) Status() HealthStatus {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.status
}

// Stop halts the background probe goroutine.
func (h *HealthCheckNotifier) Stop() {
	close(h.stop)
}

func (h *HealthCheckNotifier) loop() {
	t := time.NewTicker(h.interval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			probe := Message{Path: "__healthcheck__", Status: StatusExpiringSoon}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err := h.inner.Send(ctx, probe)
			cancel()
			h.record(err)
		case <-h.stop:
			return
		}
	}
}

func (h *HealthCheckNotifier) record(err error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.status = HealthStatus{
		Name:      h.name,
		Healthy:   err == nil,
		LastError: err,
		CheckedAt: time.Now(),
	}
}
