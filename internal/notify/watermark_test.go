package notify

import (
	"testing"
	"time"
)

func TestNewWatermarkNotifier_NilInner(t *testing.T) {
	_, err := NewWatermarkNotifier(nil, time.Hour)
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewWatermarkNotifier_ZeroDuration(t *testing.T) {
	_, err := NewWatermarkNotifier(NewNoopNotifier(), 0)
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
}

func TestNewWatermarkNotifier_NegativeDuration(t *testing.T) {
	_, err := NewWatermarkNotifier(NewNoopNotifier(), -time.Minute)
	if err == nil {
		t.Fatal("expected error for negative duration")
	}
}

func TestNewWatermarkNotifier_Valid(t *testing.T) {
	w, err := NewWatermarkNotifier(NewNoopNotifier(), time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestWatermarkNotifier_FiresOnFirstCrossing(t *testing.T) {
	var count int
	inner := &mockNotifier{fn: func(Message) error { count++; return nil }}

	w, _ := NewWatermarkNotifier(inner, time.Hour)

	msg := Message{
		Path:   "secret/db",
		Status: StatusExpiringSoon,
		Expiry: time.Now().Add(30 * time.Minute), // below 1h watermark
	}

	if err := w.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 send, got %d", count)
	}
}

func TestWatermarkNotifier_SuppressesRepeatBelowWatermark(t *testing.T) {
	var count int
	inner := &mockNotifier{fn: func(Message) error { count++; return nil }}

	w, _ := NewWatermarkNotifier(inner, time.Hour)

	msg := Message{
		Path:   "secret/db",
		Status: StatusExpiringSoon,
		Expiry: time.Now().Add(30 * time.Minute),
	}

	_ = w.Send(msg)
	_ = w.Send(msg) // second send — same path, still below
	_ = w.Send(msg) // third send

	if count != 1 {
		t.Fatalf("expected exactly 1 forward, got %d", count)
	}
}

func TestWatermarkNotifier_RefireAfterRenewal(t *testing.T) {
	var count int
	inner := &mockNotifier{fn: func(Message) error { count++; return nil }}

	w, _ := NewWatermarkNotifier(inner, time.Hour)

	below := Message{
		Path:   "secret/db",
		Status: StatusExpiringSoon,
		Expiry: time.Now().Add(30 * time.Minute),
	}
	above := Message{
		Path:   "secret/db",
		Status: StatusOK,
		Expiry: time.Now().Add(48 * time.Hour), // renewed above watermark
	}

	_ = w.Send(below) // fires (count=1)
	_ = w.Send(above) // resets state, no fire
	_ = w.Send(below) // crosses again (count=2)

	if count != 2 {
		t.Fatalf("expected 2 forwards after renewal, got %d", count)
	}
}

func TestWatermarkNotifier_AboveWatermarkSuppressed(t *testing.T) {
	var count int
	inner := &mockNotifier{fn: func(Message) error { count++; return nil }}

	w, _ := NewWatermarkNotifier(inner, time.Hour)

	msg := Message{
		Path:   "secret/db",
		Status: StatusOK,
		Expiry: time.Now().Add(72 * time.Hour), // well above watermark
	}

	_ = w.Send(msg)
	if count != 0 {
		t.Fatalf("expected no forward above watermark, got %d", count)
	}
}

func TestWatermarkNotifier_Reset(t *testing.T) {
	var count int
	inner := &mockNotifier{fn: func(Message) error { count++; return nil }}

	w, _ := NewWatermarkNotifier(inner, time.Hour)

	msg := Message{
		Path:   "secret/db",
		Status: StatusExpiringSoon,
		Expiry: time.Now().Add(20 * time.Minute),
	}

	_ = w.Send(msg) // fires
	w.Reset()       // clear state
	_ = w.Send(msg) // should fire again

	if count != 2 {
		t.Fatalf("expected 2 forwards after Reset, got %d", count)
	}
}
