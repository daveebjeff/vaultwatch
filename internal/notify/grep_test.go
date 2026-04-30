package notify

import (
	"context"
	"errors"
	"regexp"
	"testing"
)

func TestNewGrepNotifier_NilInner(t *testing.T) {
	_, err := NewGrepNotifier(nil, []*regexp.Regexp{regexp.MustCompile("foo")})
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewGrepNotifier_NoPatterns(t *testing.T) {
	_, err := NewGrepNotifier(NewNoopNotifier(), nil)
	if err == nil {
		t.Fatal("expected error for empty patterns")
	}
}

func TestNewGrepNotifier_NilPattern(t *testing.T) {
	_, err := NewGrepNotifier(NewNoopNotifier(), []*regexp.Regexp{nil})
	if err == nil {
		t.Fatal("expected error for nil pattern element")
	}
}

func TestNewGrepNotifier_Valid(t *testing.T) {
	g, err := NewGrepNotifier(NewNoopNotifier(), []*regexp.Regexp{regexp.MustCompile("ok")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestGrepNotifier_MatchForwards(t *testing.T) {
	var got Message
	cap := &capturingNotifier{fn: func(m Message) error { got = m; return nil }}
	g, _ := NewGrepNotifier(cap, []*regexp.Regexp{regexp.MustCompile(`secret`)})

	msg := Message{Path: "vault/secret", Body: "secret expires soon"}
	if err := g.Send(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Path != msg.Path {
		t.Errorf("expected message forwarded, got %+v", got)
	}
}

func TestGrepNotifier_NoMatchSuppresses(t *testing.T) {
	called := false
	cap := &capturingNotifier{fn: func(_ Message) error { called = true; return nil }}
	g, _ := NewGrepNotifier(cap, []*regexp.Regexp{regexp.MustCompile(`CRITICAL`)})

	msg := Message{Path: "vault/kv", Body: "all good"}
	if err := g.Send(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("expected message to be suppressed")
	}
}

func TestGrepNotifier_InnerErrorPropagated(t *testing.T) {
	sentinel := errors.New("downstream failure")
	cap := &capturingNotifier{fn: func(_ Message) error { return sentinel }}
	g, _ := NewGrepNotifier(cap, []*regexp.Regexp{regexp.MustCompile(`.`)})

	err := g.Send(context.Background(), Message{Body: "anything"})
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestGrepNotifier_MultiplePatterns_AnyMatch(t *testing.T) {
	var count int
	cap := &capturingNotifier{fn: func(_ Message) error { count++; return nil }}
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`alpha`),
		regexp.MustCompile(`beta`),
	}
	g, _ := NewGrepNotifier(cap, patterns)

	_ = g.Send(context.Background(), Message{Body: "beta lease expired"})
	_ = g.Send(context.Background(), Message{Body: "nothing relevant"})

	if count != 1 {
		t.Errorf("expected 1 forwarded message, got %d", count)
	}
}
