package notify

import (
	"context"
	"sync"
	"time"
)

// CooldownNotifier suppresses repeated notifications for the same secret path
// until a per-path cooldown period has elapsed. Unlike RateLimitNotifier which
// uses a fixed window, CooldownNotifier resets the timer on every suppressed
// send, ensuring a quiet period after the last event.
type CooldownNotifier struct {
	inner    Notifier
	cooldown time.Duration

	mu      sync.Mutex
	lastSent map[string]time.Time
}

// NewCooldownNotifier returns a CooldownNotifier that wraps inner and suppresses
// sends for the same path until cooldown has elapsed since the last forwarded
// message.
func NewCooldownNotifier(inner Notifier, cooldown time.Duration) (*CooldownNotifier, error) {
	if inner == nil {
		return nil, ErrNilInner
	}
	if cooldown <= 0 {
		return nil, ErrZeroCooldown
	}
	return &CooldownNotifier{
		inner:    inner,
		cooldown: cooldown,
		lastSent: make(map[string]time.Time),
	}, nil
}

// Send forwards msg to the inner notifier only if the cooldown period has
// elapsed since the last forwarded message for the same path. If the message
// is forwarded, the cooldown timer is reset.
func (c *CooldownNotifier) Send(ctx context.Context, msg Message) error {
	c.mu.Lock()
	last, seen := c.lastSent[msg.Path]
	if seen && time.Since(last) < c.cooldown {
		c.mu.Unlock()
		return nil
	}
	c.lastSent[msg.Path] = time.Now()
	c.mu.Unlock()

	return c.inner.Send(ctx, msg)
}
