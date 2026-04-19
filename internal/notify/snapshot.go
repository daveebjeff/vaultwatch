package notify

import (
	"sync"
	"time"
)

// SnapshotNotifier wraps a Notifier and records the most recent message
// sent per secret path, along with the time it was received.
type SnapshotNotifier struct {
	mu    sync.RWMutex
	inner Notifier
	snaps map[string]Snapshot
}

// Snapshot holds the last message and timestamp for a path.
type Snapshot struct {
	Message   Message
	ReceivedAt time.Time
}

// NewSnapshotNotifier wraps inner and tracks per-path snapshots.
func NewSnapshotNotifier(inner Notifier) (*SnapshotNotifier, error) {
	if inner == nil {
		return nil, ErrNilInner
	}
	return &SnapshotNotifier{
		inner: inner,
		snaps: make(map[string]Snapshot),
	}, nil
}

// Send forwards the message and records a snapshot keyed by Path.
func (s *SnapshotNotifier) Send(msg Message) error {
	err := s.inner.Send(msg)
	s.mu.Lock()
	s.snaps[msg.Path] = Snapshot{Message: msg, ReceivedAt: time.Now()}
	s.mu.Unlock()
	return err
}

// Latest returns the most recent snapshot for the given path, if any.
func (s *SnapshotNotifier) Latest(path string) (Snapshot, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snap, ok := s.snaps[path]
	return snap, ok
}

// All returns a copy of all current snapshots.
func (s *SnapshotNotifier) All() map[string]Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copy := make(map[string]Snapshot, len(s.snaps))
	for k, v := range s.snaps {
		copy[k] = v
	}
	return copy
}
