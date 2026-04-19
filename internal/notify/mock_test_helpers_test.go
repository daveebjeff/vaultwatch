package notify

// mockNotifier is a test helper that calls fn on each Send.
type mockNotifier struct {
	fn func(Message) error
}

func (m *mockNotifier) Send(msg Message) error {
	if m.fn != nil {
		return m.fn(msg)
	}
	return nil
}
