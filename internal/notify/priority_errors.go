package notify

import "errors"

// ErrNoPriorityNotifiers is returned by PriorityNotifier.Send when no
// notifiers have been registered.
var ErrNoPriorityNotifiers = errors.New("priority notifier: no notifiers registered")

// ErrNilPriorityNotifier is returned by PriorityNotifier.Add when a nil
// Notifier is provided.
var ErrNilPriorityNotifier = errors.New("priority notifier: notifier must not be nil")
