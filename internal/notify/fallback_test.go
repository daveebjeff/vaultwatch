package notify

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewFallbackNotifier_NilPrimary(t *testing.T) {
	_, err := NewFallbackNotifier(nil, NewNoopNotifier())
	if err == nil {
		t.Fatal("expected error for nil primary")
	}
}

func TestNewFallbackNotifier_NilSecondary(t *testing.T) {
	_, err := NewFallbackNotifier(NewNoopNotifier(), nil)
	if err == nil {
		t.Fatal("expected error for nil secondary")
	}
}

func TestNewFallbackNotifier_Valid(t *testing.T) {
	f, err := NewFallbackNotifier(NewNoopNotifier(), NewNoopNotifier())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestFallbackNotifier_PrimarySuccess(t *testing.T) {
	var called bool
	primary := &funcNotifier{fn: func(_ context.Context, _ Message) error { called = true; return nil }}
	secondary := &funcNotifier{fn: func(_ context.Context, _ Message) error {
		t.Fatal("secondary should not be called")
		return nil
	}}
	f, _ := NewFallbackNotifier(primary, secondary)
	if err := f.Send(context.Background(), Message{Path: "sec/key", Status: StatusExpiringSoon, ExpiresAt: time.Now().Add(time.Hour)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("primary was not called")
	}
}

func TestFallbackNotifier_PrimaryFailFallsToSecondary(t *testing.T) {
	primary := &funcNotifier{fn: func(_ context.Context, _ Message) error { return errors.New("primary down") }}
	var secondaryCalled bool
	secondary := &funcNotifier{fn: func(_ context.Context, _ Message) error { secondaryCalled = true; return nil }}
	f, _ := NewFallbackNotifier(primary, secondary)
	if err := f.Send(context.Background(), Message{Path: "sec/key", Status: StatusExpiringSoon, ExpiresAt: time.Now().Add(time.Hour)}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !secondaryCalled {
		t.Fatal("secondary was not called after primary failure")
	}
}

func TestFallbackNotifier_BothFail(t *testing.T) {
	primary := &funcNotifier{fn: func(_ context.Context, _ Message) error { return errors.New("primary down") }}
	secondary := &funcNotifier{fn: func(_ context.Context, _ Message) error { return errors.New("secondary down") }}
	f, _ := NewFallbackNotifier(primary, secondary)
	err := f.Send(context.Background(), Message{Path: "sec/key", Status: StatusExpired, ExpiresAt: time.Now()})
	if err == nil {
		t.Fatal("expected combined error")
	}
	if !errors.Is(err, errors.Unwrap(err)) && err.Error() == "" {
		t.Fatal("expected non-empty error message")
	}
}
