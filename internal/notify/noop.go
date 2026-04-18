package notify

// NoopNotifier is a Notifier that silently discards all messages.
// It is useful for testing and as a safe default when no notifier is configured.
type NoopNotifier struct{}

// NewNoopNotifier returns a NoopNotifier.
func NewNoopNotifier() *NoopNotifier {
	return &NoopNotifier{}
}

// Send discards the message and returns nil.
func (n *NoopNotifier) Send(msg Message) error {
	return nil
}
