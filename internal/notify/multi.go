package notify

import "fmt"

// MultiNotifier fans a single message out to multiple Notifier backends.
type MultiNotifier struct {
	notifiers []Notifier
}

// NewMultiNotifier creates a MultiNotifier from the provided backends.
func NewMultiNotifier(ns ...Notifier) *MultiNotifier {
	return &MultiNotifier{notifiers: ns}
}

// Add appends a new backend at runtime.
func (m *MultiNotifier) Add(n Notifier) {
	m.notifiers = append(m.notifiers, n)
}

// Send delivers msg to every registered backend, collecting errors.
func (m *MultiNotifier) Send(msg Message) error {
	var errs []error
	for _, n := range m.notifiers {
		if err := n.Send(msg); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("notify: %d backend(s) failed: %v", len(errs), errs)
}
