package notify

import (
	"fmt"
	"sort"
	"sync"
)

// PriorityNotifier routes messages to one of several notifiers based on a
// priority level derived from the message status. Higher-priority notifiers
// are tried first; if a higher-priority notifier fails the error is returned
// without falling back.
type PriorityNotifier struct {
	mu      sync.RWMutex
	buckets []priorityBucket // sorted ascending by level (lowest first)
}

type priorityBucket struct {
	level    int
	notifier Notifier
}

// NewPriorityNotifier returns an empty PriorityNotifier. Use Add to register
// notifiers at specific priority levels. Lower level numbers are lower
// priority; the notifier with the highest level whose threshold is met is
// used.
func NewPriorityNotifier() *PriorityNotifier {
	return &PriorityNotifier{}
}

// Add registers a Notifier at the given priority level. If two notifiers share
// the same level the one added last wins.
func (p *PriorityNotifier) Add(level int, n Notifier) error {
	if n == nil {
		return fmt.Errorf("priority notifier: notifier at level %d must not be nil", level)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, b := range p.buckets {
		if b.level == level {
			p.buckets[i].notifier = n
			return nil
		}
	}
	p.buckets = append(p.buckets, priorityBucket{level: level, notifier: n})
	sort.Slice(p.buckets, func(i, j int) bool {
		return p.buckets[i].level < p.buckets[j].level
	})
	return nil
}

// Send dispatches msg to the highest-priority registered notifier whose level
// is greater than or equal to the numeric status severity (Expired=2,
// ExpiringSoon=1, OK=0). If no bucket matches the lowest-level notifier is
// used as a fallback.
func (p *PriorityNotifier) Send(msg Message) error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if len(p.buckets) == 0 {
		return fmt.Errorf("priority notifier: no notifiers registered")
	}
	threshold := statusSeverity(msg.Status)
	// Walk from highest to lowest, pick first whose level >= threshold.
	for i := len(p.buckets) - 1; i >= 0; i-- {
		if p.buckets[i].level >= threshold {
			return p.buckets[i].notifier.Send(msg)
		}
	}
	// Fallback: lowest bucket.
	return p.buckets[0].notifier.Send(msg)
}

func statusSeverity(s Status) int {
	switch s {
	case StatusExpired:
		return 2
	case StatusExpiringSoon:
		return 1
	default:
		return 0
	}
}
