package notify

import (
	"context"
	"errors"
	"testing"
	"time"
)

var seqMsg = Message{
	Path:      "secret/data/db",
	Status:    StatusExpiringSoon,
	ExpiresAt: time.Now().Add(10 * time.Minute),
}

func TestNewSequenceNotifier_NoSteps(t *testing.T) {
	_, err := NewSequenceNotifier()
	if err == nil {
		t.Fatal("expected error for empty steps")
	}
}

func TestNewSequenceNotifier_NilStep(t *testing.T) {
	_, err := NewSequenceNotifier(NewNoopNotifier(), nil)
	if err == nil {
		t.Fatal("expected error for nil step")
	}
}

func TestNewSequenceNotifier_Valid(t *testing.T) {
	sn, err := NewSequenceNotifier(NewNoopNotifier(), NewNoopNotifier())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sn.Len() != 2 {
		t.Fatalf("expected 2 steps, got %d", sn.Len())
	}
}

func TestSequenceNotifier_AllCalledOnSuccess(t *testing.T) {
	var calls []int
	make := func(id int) Notifier {
		return &mockNotifier{fn: func(_ context.Context, _ Message) error {
			calls = append(calls, id)
			return nil
		}}
	}
	sn, _ := NewSequenceNotifier(make(1), make(2), make(3))
	if err := sn.Send(context.Background(), seqMsg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(calls) != 3 {
		t.Fatalf("expected 3 calls, got %d", len(calls))
	}
	for i, v := range calls {
		if v != i+1 {
			t.Errorf("step %d: expected id %d, got %d", i, i+1, v)
		}
	}
}

func TestSequenceNotifier_StopsOnFirstError(t *testing.T) {
	var calls int
	boom := errors.New("boom")
	step1 := &mockNotifier{fn: func(_ context.Context, _ Message) error { calls++; return nil }}
	step2 := &mockNotifier{fn: func(_ context.Context, _ Message) error { calls++; return boom }}
	step3 := &mockNotifier{fn: func(_ context.Context, _ Message) error { calls++; return nil }}

	sn, _ := NewSequenceNotifier(step1, step2, step3)
	err := sn.Send(context.Background(), seqMsg)
	if err == nil {
		t.Fatal("expected error from step2")
	}
	if !errors.Is(err, boom) {
		t.Errorf("expected wrapped boom error, got: %v", err)
	}
	if calls != 2 {
		t.Errorf("expected 2 calls (step3 must not run), got %d", calls)
	}
}

func TestSequenceNotifier_ErrorIncludesStepIndex(t *testing.T) {
	fail := &mockNotifier{fn: func(_ context.Context, _ Message) error {
		return errors.New("fail")
	}}
	sn, _ := NewSequenceNotifier(NewNoopNotifier(), NewNoopNotifier(), fail)
	err := sn.Send(context.Background(), seqMsg)
	if err == nil {
		t.Fatal("expected error")
	}
	const want = "sequence: step 2 failed"
	if !containsSubstring(err.Error(), want) {
		t.Errorf("error %q does not contain %q", err.Error(), want)
	}
}

func containsSubstring(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
