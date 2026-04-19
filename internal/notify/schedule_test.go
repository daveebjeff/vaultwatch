package notify

import (
	"testing"
	"time"
)

func TestNewScheduleNotifier_NilInner(t *testing.T) {
	_, err := NewScheduleNotifier(nil, nil, TimeWindow{Start: 0, End: 8 * time.Hour})
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewScheduleNotifier_NoWindows(t *testing.T) {
	n := NewNoopNotifier()
	_, err := NewScheduleNotifier(n, nil)
	if err == nil {
		t.Fatal("expected error for no windows")
	}
}

func TestNewScheduleNotifier_Valid(t *testing.T) {
	n := NewNoopNotifier()
	s, err := NewScheduleNotifier(n, time.UTC, TimeWindow{Start: 0, End: 24 * time.Hour})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestScheduleNotifier_AlwaysOpen(t *testing.T) {
	var called bool
	n := &mockNotifier{fn: func(msg Message) error { called = true; return nil }}
	s, _ := NewScheduleNotifier(n, time.UTC, TimeWindow{Start: 0, End: 24 * time.Hour})
	if err := s.Send(Message{Path: "secret/a", Status: StatusExpired}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected inner to be called")
	}
}

func TestScheduleNotifier_AlwaysClosed(t *testing.T) {
	var called bool
	n := &mockNotifier{fn: func(msg Message) error { called = true; return nil }}
	// window in the past relative to any time: 0 duration window
	s, _ := NewScheduleNotifier(n, time.UTC, TimeWindow{Start: 25 * time.Hour, End: 26 * time.Hour})
	if err := s.Send(Message{Path: "secret/a", Status: StatusExpired}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Fatal("expected inner NOT to be called outside window")
	}
}

func TestScheduleNotifier_NilLocation_UsesUTC(t *testing.T) {
	n := NewNoopNotifier()
	s, err := NewScheduleNotifier(n, nil, TimeWindow{Start: 0, End: 24 * time.Hour})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.location != time.UTC {
		t.Fatal("expected UTC location")
	}
}
