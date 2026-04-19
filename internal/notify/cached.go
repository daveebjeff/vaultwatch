package notify

import (
	"sync"
	"time"
)

// CachedNotifier wraps a Notifier and caches the last successful send
// result for a configurable TTL. Duplicate messages within the TTL
// are forwarded only once.
type CachedNotifier struct {
	inner    Notifier
	ttl      time.Duration
	mu       sync.Mutex
	cache    map[string]time.Time
}

// NewCachedNotifier returns a CachedNotifier that suppresses duplicate
// messages (same path + status) within ttl duration.
func NewCachedNotifier(inner Notifier, ttl time.Duration) (*CachedNotifier, error) {
	if inner == nil {
		return nil, ErrNilInner
	}
	if ttl <= 0 {
		return nil, ErrZeroCooldown
	}
	return &CachedNotifier{
		inner: inner,
		ttl:   ttl,
		cache: make(map[string]time.Time),
	}, nil
}

// Send forwards msg to the inner notifier only if the same path+status
// combination has not been sent within the configured TTL.
func (c *CachedNotifier) Send(msg Message) error {
	key := msg.Path + "|" + string(msg.Status)

	c.mu.Lock()
	last, ok := c.cache[key]
	if ok && time.Since(last) < c.ttl {
		c.mu.Unlock()
		return nil
	}
	c.cache[key] = time.Now()
	c.mu.Unlock()

	return c.inner.Send(msg)
}

// Invalidate removes a cached entry for the given path and status,
// allowing the next send to be forwarded regardless of TTL.
func (c *CachedNotifier) Invalidate(path string, status Status) {
	key := path + "|" + string(status)
	c.mu.Lock()
	delete(c.cache, key)
	c.mu.Unlock()
}
