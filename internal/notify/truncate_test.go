package notify

import (
	"errors"
	"strings"
	"testing"
)

func TestNewTruncateNotifier_NilInner(t *testing.T) {
	_, err := NewTruncateNotifier(nil, 100, "…")
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewTruncateNotifier_ZeroMaxRunes(t *testing.T) {
	noop := NewNoopNotifier()
	_, err := NewTruncateNotifier(noop, 0, "")
	if err == nil {
		t.Fatal("expected error for maxRunes=0")
	}
}

func TestNewTruncateNotifier_SuffixTooLong(t *testing.T) {
	noop := NewNoopNotifier()
	_, err := NewTruncateNotifier(noop, 3, "toolong")
	if err == nil {
		t.Fatal("expected error when suffix >= maxRunes")
	}
}

func TestNewTruncateNotifier_Valid(t *testing.T) {
	noop := NewNoopNotifier()
	tn, err := NewTruncateNotifier(noop, 80, "…")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tn == nil {
		t.Fatal("expected non-nil TruncateNotifier")
	}
}

func TestTruncateNotifier_ShortBodyUnchanged(t *testing.T) {
	cap := &captureNotifier{}
	tn, _ := NewTruncateNotifier(cap, 100, "…")
	msg := Message{Path: "secret/a", Body: "short"}
	if err := tn.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cap.last.Body != "short" {
		t.Errorf("expected body unchanged, got %q", cap.last.Body)
	}
}

func TestTruncateNotifier_LongBodyTruncated(t *testing.T) {
	cap := &captureNotifier{}
	tn, _ := NewTruncateNotifier(cap, 10, "…")
	msg := Message{Path: "secret/b", Body: strings.Repeat("a", 20)}
	if err := tn.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := cap.last.Body
	if len([]rune(got)) != 10 {
		t.Errorf("expected 10 runes, got %d: %q", len([]rune(got)), got)
	}
	if !strings.HasSuffix(got, "…") {
		t.Errorf("expected suffix '…', got %q", got)
	}
}

func TestTruncateNotifier_InnerErrorPropagated(t *testing.T) {
	fail := &failNotifier{err: errors.New("send failed")}
	tn, _ := NewTruncateNotifier(fail, 50, "")
	err := tn.Send(Message{Path: "secret/c", Body: "hello"})
	if err == nil {
		t.Fatal("expected error from inner notifier")
	}
}

func TestTruncateNotifier_ExactLengthUnchanged(t *testing.T) {
	cap := &captureNotifier{}
	tn, _ := NewTruncateNotifier(cap, 5, "…")
	msg := Message{Path: "secret/d", Body: "hello"}
	if err := tn.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cap.last.Body != "hello" {
		t.Errorf("expected body unchanged at exact limit, got %q", cap.last.Body)
	}
}
