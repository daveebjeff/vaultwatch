package notify

import (
	"fmt"
	"sync"
	"time"
)

// ScheduleNotifier suppresses notifications outside of allowed time windows.
// For example, only forward alerts during business hours.
type ScheduleNotifier struct {
	inner    Notifier
	windows  []TimeWindow
	location *time.Location
	mu       sync.Mutex
}

// TimeWindow defines a daily time range using hour:minute boundaries.
type TimeWindow struct {
	Start time.Duration // offset from midnight
	End   time.Duration // offset from midnight
}

// NewScheduleNotifier creates a ScheduleNotifier that forwards only within the given windows.
// loc may be nil, in which case UTC is used.
func NewScheduleNotifier(inner Notifier, loc *time.Location, windows ...TimeWindow) (*ScheduleNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("schedule: inner notifier must not be nil")
	}
	if len(windows) == 0 {
		return nil, fmt.Errorf("schedule: at least one time window is required")
	}
	if loc == nil {
		loc = time.UTC
	}
	return &ScheduleNotifier{inner: inner, windows: windows, location: loc}, nil
}

// Send forwards the message only if the current time falls within an allowed window.
func (s *ScheduleNotifier) Send(msg Message) error {
	s.mu.Lock()
	loc := s.location
	windows := s.windows
	s.mu.Unlock()

	now := time.Now().In(loc)
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	offset := now.Sub(midnight)

	for _, w := range windows {
		if offset >= w.Start && offset < w.End {
			return s.inner.Send(msg)
		}
	}
	return nil
}
