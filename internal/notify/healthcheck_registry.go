package notify

import (
	"fmt"
	"sync"
)

// HealthRegistry tracks multiple HealthCheckNotifiers by name.
type HealthRegistry struct {
	mu      sync.RWMutex
	entries map[string]*HealthCheckNotifier
}

// NewHealthRegistry returns an empty HealthRegistry.
func NewHealthRegistry() *HealthRegistry {
	return &HealthRegistry{entries: make(map[string]*HealthCheckNotifier)}
}

// Register adds a HealthCheckNotifier under its name.
// Returns an error if the name is already registered.
func (r *HealthRegistry) Register(h *HealthCheckNotifier) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.entries[h.name]; exists {
		return fmt.Errorf("healthcheck: notifier %q already registered", h.name)
	}
	r.entries[h.name] = h
	return nil
}

// Statuses returns a snapshot of all registered health statuses.
func (r *HealthRegistry) Statuses() []HealthStatus {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]HealthStatus, 0, len(r.entries))
	for _, h := range r.entries {
		out = append(out, h.Status())
	}
	return out
}

// AllHealthy returns true only if every registered notifier is healthy.
func (r *HealthRegistry) AllHealthy() bool {
	for _, s := range r.Statuses() {
		if !s.Healthy {
			return false
		}
	}
	return true
}

// StopAll stops the background probe goroutine for every registered notifier.
func (r *HealthRegistry) StopAll() {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, h := range r.entries {
		h.Stop()
	}
}
