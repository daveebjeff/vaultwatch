package notify

import (
	"context"
	"sync"
	"time"
)

// MuteNotifier wraps a Notifier and suppresses all sends during a mute window.
// Muting can be activated for a fixed duration or toggled manually.
// Once the mute window expires, sends are forwarded again automatically.
type MuteNotifier struct {
	inner    Notifier
	mu       sync.Mutex
	muteUntil time.Time
}

// NewMuteNotifier creates a MuteNotifier wrapping inner.
// Returns an error if inner is nil.
func NewMuteNotifier(inner Notifier) (*MuteNotifier, error) {
	if inner == nil {
		return nil, ErrNilInner
	}
	return &MuteNotifier{inner: inner}, nil
}

// Mute silences all notifications for the given duration.
// Calling Mute while already muted extends the mute window to
// whichever deadline is later.
func (m *MuteNotifier) Mute(d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	candidate := time.Now().Add(d)
	if candidate.After(m.muteUntil) {
		m.muteUntil = candidate
	}
}

// Unmute cancels any active mute window immediately.
func (m *MuteNotifier) Unmute() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.muteUntil = time.Time{}
}

// IsMuted reports whether the notifier is currently muted.
func (m *MuteNotifier) IsMuted() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return time.Now().Before(m.muteUntil)
}

// Send forwards msg to the inner notifier unless a mute window is active,
// in which case the message is silently dropped and nil is returned.
func (m *MuteNotifier) Send(ctx context.Context, msg Message) error {
	if m.IsMuted() {
		return nil
	}
	return m.inner.Send(ctx, msg)
}
