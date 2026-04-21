package notify

import "context"

// TeeNotifier forwards every message to two notifiers simultaneously,
// similar to the Unix tee command. Both notifiers always receive the
// message regardless of whether the other returns an error. If both
// return errors they are combined via errors.Join.
type TeeNotifier struct {
	a Notifier
	b Notifier
}

// NewTeeNotifier returns a TeeNotifier that sends each message to both
// a and b. Neither notifier may be nil.
func NewTeeNotifier(a, b Notifier) (*TeeNotifier, error) {
	if a == nil {
		return nil, ErrNilNotifier
	}
	if b == nil {
		return nil, ErrNilNotifier
	}
	return &TeeNotifier{a: a, b: b}, nil
}

// Send delivers msg to both notifiers. Both are always called. If
// one or both return an error, the errors are joined and returned.
func (t *TeeNotifier) Send(ctx context.Context, msg Message) error {
	errA := t.a.Send(ctx, msg)
	errB := t.b.Send(ctx, msg)
	if errA != nil && errB != nil {
		return joinErrors(errA, errB)
	}
	if errA != nil {
		return errA
	}
	return errB
}
