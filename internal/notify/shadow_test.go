package notify

import (
	"context"
	"errors"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

func TestNewShadowNotifier_NilPrimary(t *testing.T) {
	_, err := NewShadowNotifier(nil, NewNoopNotifier(), nil)
	if !errors.Is(err, ErrNilNotifier) {
		t.Fatalf("expected ErrNilNotifier, got %v", err)
	}
}

func TestNewShadowNotifier_NilShadow(t *testing.T) {
	_, err := NewShadowNotifier(NewNoopNotifier(), nil, nil)
	if !errors.Is(err, ErrNilNotifier) {
		t.Fatalf("expected ErrNilNotifier, got %v", err)
	}
}

func TestNewShadowNotifier_Valid(t *testing.T) {
	s, err := NewShadowNotifier(NewNoopNotifier(), NewNoopNotifier(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil ShadowNotifier")
	}
}

func TestShadowNotifier_PrimaryErrorReturned(t *testing.T) {
	primaryErr := errors.New("primary failure")
	primary := &callbackNotifier{fn: func(_ context.Context, _ Message) error { return primaryErr }}
	shadow := NewNoopNotifier()

	s, _ := NewShadowNotifier(primary, shadow, log.New(os.Stderr, "", 0))
	err := s.Send(context.Background(), Message{Path: "secret/test"})
	if !errors.Is(err, primaryErr) {
		t.Fatalf("expected primaryErr, got %v", err)
	}
}

func TestShadowNotifier_ShadowErrorSuppressed(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	shadow := &callbackNotifier{fn: func(_ context.Context, _ Message) error {
		defer wg.Done()
		return errors.New("shadow failure")
	}}

	s, _ := NewShadowNotifier(NewNoopNotifier(), shadow, log.New(os.Stderr, "", 0))
	err := s.Send(context.Background(), Message{Path: "secret/test"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("shadow goroutine did not complete in time")
	}
}

func TestShadowNotifier_BothCalled(t *testing.T) {
	var mu sync.Mutex
	called := map[string]bool{}

	make := func(name string) *callbackNotifier {
		return &callbackNotifier{fn: func(_ context.Context, _ Message) error {
			mu.Lock(); called[name] = true; mu.Unlock()
			return nil
		}}
	}

	s, _ := NewShadowNotifier(make("primary"), make("shadow"), nil)
	s.Send(context.Background(), Message{Path: "secret/x"})
	time.Sleep(50 * time.Millisecond)

	mu.Lock(); defer mu.Unlock()
	if !called["primary"] || !called["shadow"] {
		t.Fatalf("expected both called, got %v", called)
	}
}
