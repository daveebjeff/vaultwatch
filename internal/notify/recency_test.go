package notify

import (
	"errors"
	"testing"
	"time"
)

func TestNewRecencyNotifier_NilInner(t *testing.T) {
	_, err := NewRecencyNotifier(nil, time.Minute)
	if err == nil {
		t.Fatal("expected error for nil inner, got nil")
	}
}

func TestNewRecencyNotifier_ZeroWindow(t *testing.T) {
	_, err := NewRecencyNotifier(NewNoopNotifier(), 0)
	if err == nil {
		t.Fatal("expected error for zero window, got nil")
	}
}

func TestNewRecencyNotifier_Valid(t *testing.T) {
	r, err := NewRecencyNotifier(NewNoopNotifier(), time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestRecencyNotifier_FirstSendForwarded(t *testing.T) {
	mock := &mockNotifier{}
	r, _ := NewRecencyNotifier(mock, time.Minute)

	msg := Message{Path: "secret/foo", Status: StatusExpiringSoon}
	if err := r.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.calls != 1 {
		t.Fatalf("expected 1 call, got %d", mock.calls)
	}
}

func TestRecencyNotifier_DuplicateSuppressed(t *testing.T) {
	mock := &mockNotifier{}
	r, _ := NewRecencyNotifier(mock, time.Minute)

	msg := Message{Path: "secret/foo", Status: StatusExpiringSoon}
	_ = r.Send(msg)
	_ = r.Send(msg)

	if mock.calls != 1 {
		t.Fatalf("expected 1 call after duplicate, got %d", mock.calls)
	}
}

func TestRecencyNotifier_StatusChangeForwarded(t *testing.T) {
	mock := &mockNotifier{}
	r, _ := NewRecencyNotifier(mock, time.Minute)

	_ = r.Send(Message{Path: "secret/foo", Status: StatusExpiringSoon})
	_ = r.Send(Message{Path: "secret/foo", Status: StatusExpired})

	if mock.calls != 2 {
		t.Fatalf("expected 2 calls on status change, got %d", mock.calls)
	}
}

func TestRecencyNotifier_WindowExpiredForwards(t *testing.T) {
	mock := &mockNotifier{}
	r, _ := NewRecencyNotifier(mock, 10*time.Millisecond)

	msg := Message{Path: "secret/bar", Status: StatusExpiringSoon}
	_ = r.Send(msg)
	time.Sleep(20 * time.Millisecond)
	_ = r.Send(msg)

	if mock.calls != 2 {
		t.Fatalf("expected 2 calls after window expiry, got %d", mock.calls)
	}
}

func TestRecencyNotifier_Reset(t *testing.T) {
	mock := &mockNotifier{}
	r, _ := NewRecencyNotifier(mock, time.Minute)

	msg := Message{Path: "secret/baz", Status: StatusExpiringSoon}
	_ = r.Send(msg)
	r.Reset()
	_ = r.Send(msg)

	if mock.calls != 2 {
		t.Fatalf("expected 2 calls after Reset, got %d", mock.calls)
	}
}

func TestRecencyNotifier_InnerErrorReturned(t *testing.T) {
	expected := errors.New("send failed")
	mock := &mockNotifier{err: expected}
	r, _ := NewRecencyNotifier(mock, time.Minute)

	err := r.Send(Message{Path: "secret/err", Status: StatusExpired})
	if !errors.Is(err, expected) {
		t.Fatalf("expected inner error, got %v", err)
	}
}
