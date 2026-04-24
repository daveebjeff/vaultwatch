package notify

import "fmt"

// PriorityEntry describes a single level-to-notifier binding used when
// constructing a PriorityNotifier from configuration.
type PriorityEntry struct {
	// Level is the numeric priority. Higher values are higher priority.
	Level int
	// Notifier is the Notifier to invoke when this level is selected.
	Notifier Notifier
}

// BuildPriorityNotifier constructs a PriorityNotifier from a slice of
// PriorityEntry values. It returns an error if any entry contains a nil
// Notifier or if no entries are provided.
func BuildPriorityNotifier(entries []PriorityEntry) (*PriorityNotifier, error) {
	if len(entries) == 0 {
		return nil, fmt.Errorf("priority notifier: at least one entry is required")
	}
	p := NewPriorityNotifier()
	for _, e := range entries {
		if err := p.Add(e.Level, e.Notifier); err != nil {
			return nil, err
		}
	}
	return p, nil
}
