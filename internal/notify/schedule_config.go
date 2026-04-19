package notify

import (
	"fmt"
	"time"
)

// ScheduleConfig holds the serialisable form of a ScheduleNotifier configuration.
type ScheduleConfig struct {
	Timezone string              `yaml:"timezone" json:"timezone"`
	Windows  []TimeWindowConfig  `yaml:"windows"  json:"windows"`
}

// TimeWindowConfig is the serialisable form of a TimeWindow.
type TimeWindowConfig struct {
	Start string `yaml:"start" json:"start"` // "HH:MM"
	End   string `yaml:"end"   json:"end"`   // "HH:MM"
}

// Build converts a ScheduleConfig into a ScheduleNotifier wrapping inner.
func (c ScheduleConfig) Build(inner Notifier) (*ScheduleNotifier, error) {
	loc := time.UTC
	if c.Timezone != "" {
		var err error
		loc, err = time.LoadLocation(c.Timezone)
		if err != nil {
			return nil, fmt.Errorf("schedule config: invalid timezone %q: %w", c.Timezone, err)
		}
	}
	windows := make([]TimeWindow, 0, len(c.Windows))
	for _, wc := range c.Windows {
		start, err := parseDayOffset(wc.Start)
		if err != nil {
			return nil, fmt.Errorf("schedule config: invalid start %q: %w", wc.Start, err)
		}
		end, err := parseDayOffset(wc.End)
		if err != nil {
			return nil, fmt.Errorf("schedule config: invalid end %q: %w", wc.End, err)
		}
		windows = append(windows, TimeWindow{Start: start, End: end})
	}
	return NewScheduleNotifier(inner, loc, windows...)
}

func parseDayOffset(s string) (time.Duration, error) {
	t, err := time.Parse("15:04", s)
	if err != nil {
		return 0, err
	}
	return time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute, nil
}
