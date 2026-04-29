package notify

import (
	"context"
	"sync"
	"time"
)

// CheckpointNotifier wraps an inner Notifier and records the last successful
// send time and message for each secret path. It exposes LastSeen so that
// other components (e.g. a status page) can query per-path health.
type CheckpointNotifier struct {
	inner   Notifier
	mu      sync.RWMutex
	records map[string]CheckpointRecord
}

// CheckpointRecord holds the last observed state for a single secret path.
type CheckpointRecord struct {
	Path      string
	Status    Status
	SentAt    time.Time
	Succeeded bool
}

// NewCheckpointNotifier wraps inner and begins tracking per-path checkpoints.
// It returns an error if inner is nil.
func NewCheckpointNotifier(inner Notifier) (*CheckpointNotifier, error) {
	if inner == nil {
		return nil, ErrNilInner
	}
	return &CheckpointNotifier{
		inner:   inner,
		records: make(map[string]CheckpointRecord),
	}, nil
}

// Send forwards msg to the inner notifier and records the outcome.
func (c *CheckpointNotifier) Send(ctx context.Context, msg Message) error {
	err := c.inner.Send(ctx, msg)

	c.mu.Lock()
	c.records[msg.Path] = CheckpointRecord{
		Path:      msg.Path,
		Status:    msg.Status,
		SentAt:    time.Now(),
		Succeeded: err == nil,
	}
	c.mu.Unlock()

	return err
}

// LastSeen returns the most recent CheckpointRecord for path and whether one
// exists.
func (c *CheckpointNotifier) LastSeen(path string) (CheckpointRecord, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	r, ok := c.records[path]
	return r, ok
}

// All returns a snapshot of every tracked CheckpointRecord.
func (c *CheckpointNotifier) All() []CheckpointRecord {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]CheckpointRecord, 0, len(c.records))
	for _, r := range c.records {
		out = append(out, r)
	}
	return out
}
