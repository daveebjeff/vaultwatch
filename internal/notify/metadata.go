package notify

import (
	"fmt"
	"sync"
)

// MetadataNotifier wraps an inner Notifier and attaches static key-value
// metadata to every outgoing Message before forwarding it.
//
// Metadata keys are merged with any labels already present on the message;
// existing keys are NOT overwritten, giving per-message labels priority.
type MetadataNotifier struct {
	inner    Notifier
	metadata map[string]string
	mu       sync.RWMutex
}

// NewMetadataNotifier returns a MetadataNotifier that stamps each message
// with the provided metadata before delegating to inner.
func NewMetadataNotifier(inner Notifier, metadata map[string]string) (*MetadataNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("metadata notifier: inner notifier must not be nil")
	}
	if len(metadata) == 0 {
		return nil, fmt.Errorf("metadata notifier: at least one metadata entry is required")
	}
	copy := make(map[string]string, len(metadata))
	for k, v := range metadata {
		copy[k] = v
	}
	return &MetadataNotifier{inner: inner, metadata: copy}, nil
}

// Send stamps msg.Labels with the configured metadata (without overwriting
// existing keys) and forwards the enriched message to the inner notifier.
func (m *MetadataNotifier) Send(msg Message) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if msg.Labels == nil {
		msg.Labels = make(map[string]string, len(m.metadata))
	}
	for k, v := range m.metadata {
		if _, exists := msg.Labels[k]; !exists {
			msg.Labels[k] = v
		}
	}
	return m.inner.Send(msg)
}

// SetMetadata atomically replaces the metadata map.
func (m *MetadataNotifier) SetMetadata(metadata map[string]string) error {
	if len(metadata) == 0 {
		return fmt.Errorf("metadata notifier: replacement metadata must not be empty")
	}
	copy := make(map[string]string, len(metadata))
	for k, v := range metadata {
		copy[k] = v
	}
	m.mu.Lock()
	m.metadata = copy
	m.mu.Unlock()
	return nil
}
