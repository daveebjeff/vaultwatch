package notify

import (
	"context"
	"regexp"
	"testing"
)

func TestNewRedactNotifier_NilInner(t *testing.T) {
	_, err := NewRedactNotifier(nil, MustCompilePatterns([]string{`foo`}), "")
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewRedactNotifier_NoPatterns(t *testing.T) {
	_, err := NewRedactNotifier(NewNoopNotifier(), []*regexp.Regexp{}, "")
	if err == nil {
		t.Fatal("expected error for empty patterns")
	}
}

func TestNewRedactNotifier_DefaultReplacement(t *testing.T) {
	n, err := NewRedactNotifier(NewNoopNotifier(), MustCompilePatterns([]string{`secret`}), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.replacement != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", n.replacement)
	}
}

func TestNewRedactNotifier_CustomReplacement(t *testing.T) {
	n, err := NewRedactNotifier(NewNoopNotifier(), MustCompilePatterns([]string{`secret`}), "***")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.replacement != "***" {
		t.Errorf("expected ***, got %q", n.replacement)
	}
}

func TestRedactNotifier_Send_RedactsBody(t *testing.T) {
	var got Message
	tap, _ := NewTapNotifier(NewNoopNotifier(), func(_ context.Context, m Message) {
		got = m
	})

	patterns := MustCompilePatterns([]string{`(?i)token=\S+`})
	n, err := NewRedactNotifier(tap, patterns, "[REDACTED]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msg := Message{Path: "secret/db", Body: "token=s.abc123 renewed"}
	if err := n.Send(context.Background(), msg); err != nil {
		t.Fatalf("Send error: %v", err)
	}

	if got.Body == msg.Body {
		t.Error("expected body to be redacted")
	}
	expected := "[REDACTED] renewed"
	if got.Body != expected {
		t.Errorf("got body %q, want %q", got.Body, expected)
	}
}

func TestRedactNotifier_Send_PathUnchanged(t *testing.T) {
	var got Message
	tap, _ := NewTapNotifier(NewNoopNotifier(), func(_ context.Context, m Message) {
		got = m
	})

	patterns := MustCompilePatterns([]string{`secret`})
	n, _ := NewRedactNotifier(tap, patterns, "[REDACTED]")

	msg := Message{Path: "secret/db", Body: "all fine"}
	_ = n.Send(context.Background(), msg)

	if got.Path != "secret/db" {
		t.Errorf("path should be unchanged, got %q", got.Path)
	}
}

func TestCompilePatterns_InvalidRegex(t *testing.T) {
	_, err := CompilePatterns([]string{`[invalid`})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNewDefaultRedactNotifier_Valid(t *testing.T) {
	n, err := NewDefaultRedactNotifier(NewNoopNotifier())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestDefaultRedactNotifier_RedactsVaultToken(t *testing.T) {
	var got Message
	tap, _ := NewTapNotifier(NewNoopNotifier(), func(_ context.Context, m Message) {
		got = m
	})

	n, _ := NewDefaultRedactNotifier(tap)
	msg := Message{
		Path: "auth/token",
		Body: "lease renewed for s.ABCDEFGHIJKLMNOPQRSTUVWX",
	}
	_ = n.Send(context.Background(), msg)

	if got.Body == msg.Body {
		t.Errorf("expected Vault token to be redacted; body=%q", got.Body)
	}
}
