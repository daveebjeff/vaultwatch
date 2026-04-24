package notify

import (
	"errors"
	"testing"
	"time"
)

func priorityMsg(s Status) Message {
	return Message{
		Path:      "secret/test",
		Status:    s,
		ExpiresAt: time.Now().Add(time.Hour),
	}
}

func TestNewPriorityNotifier_Empty(t *testing.T) {
	p := NewPriorityNotifier()
	if p == nil {
		t.Fatal("expected non-nil PriorityNotifier")
	}
}

func TestPriorityNotifier_Add_NilReturnsError(t *testing.T) {
	p := NewPriorityNotifier()
	if err := p.Add(1, nil); err == nil {
		t.Fatal("expected error for nil notifier")
	}
}

func TestPriorityNotifier_Send_NoNotifiers(t *testing.T) {
	p := NewPriorityNotifier()
	err := p.Send(priorityMsg(StatusExpired))
	if err == nil {
		t.Fatal("expected error with no notifiers registered")
	}
}

func TestPriorityNotifier_Send_RoutesToHighest(t *testing.T) {
	low := &mockNotifier{}
	high := &mockNotifier{}
	p := NewPriorityNotifier()
	_ = p.Add(1, low)
	_ = p.Add(2, high)

	// Expired has severity 2, should route to high (level 2).
	_ = p.Send(priorityMsg(StatusExpired))
	if high.calls != 1 {
		t.Fatalf("expected high notifier called once, got %d", high.calls)
	}
	if low.calls != 0 {
		t.Fatalf("expected low notifier not called, got %d", low.calls)
	}
}

func TestPriorityNotifier_Send_FallsBackToLowest(t *testing.T) {
	low := &mockNotifier{}
	p := NewPriorityNotifier()
	_ = p.Add(1, low)

	// Only one bucket, always used.
	_ = p.Send(priorityMsg(StatusOK))
	if low.calls != 1 {
		t.Fatalf("expected low notifier called once, got %d", low.calls)
	}
}

func TestPriorityNotifier_Send_ReturnsInnerError(t *testing.T) {
	sentinel := errors.New("notifier failure")
	n := &mockNotifier{err: sentinel}
	p := NewPriorityNotifier()
	_ = p.Add(2, n)

	err := p.Send(priorityMsg(StatusExpired))
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestPriorityNotifier_Add_SameLevelReplaces(t *testing.T) {
	first := &mockNotifier{}
	second := &mockNotifier{}
	p := NewPriorityNotifier()
	_ = p.Add(1, first)
	_ = p.Add(1, second)

	_ = p.Send(priorityMsg(StatusExpired))
	if second.calls != 1 {
		t.Fatalf("expected second notifier called, got %d", second.calls)
	}
	if first.calls != 0 {
		t.Fatalf("expected first notifier not called, got %d", first.calls)
	}
}

func TestStatusSeverity(t *testing.T) {
	tests := []struct {
		status   Status
		wantSev  int
	}{
		{StatusExpired, 2},
		{StatusExpiringSoon, 1},
		{StatusOK, 0},
	}
	for _, tc := range tests {
		got := statusSeverity(tc.status)
		if got != tc.wantSev {
			t.Errorf("statusSeverity(%v) = %d, want %d", tc.status, got, tc.wantSev)
		}
	}
}
