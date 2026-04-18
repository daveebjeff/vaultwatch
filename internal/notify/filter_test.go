package notify

import (
	"errors"
	"testing"
	"time"
)

type captureNotifier struct {
	msgs []Message
	err  error
}

func (c *captureNotifier) Send(msg Message) error {
	c.msgs = append(c.msgs, msg)
	return c.err
}

func baseMsg(path string) Message {
	return Message{
		Path:      path,
		Status:    StatusExpiringSoon,
		ExpiresAt: time.Now().Add(time.Hour),
	}
}

func TestNewFilterNotifier_NilInner(t *testing.T) {
	_, err := NewFilterNotifier(nil, []string{"secret/"})
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewFilterNotifier_NoPrefixes(t *testing.T) {
	_, err := NewFilterNotifier(&captureNotifier{}, nil)
	if err == nil {
		t.Fatal("expected error for empty prefixes")
	}
}

func TestFilterNotifier_MatchForwards(t *testing.T) {
	cap := &captureNotifier{}
	f, _ := NewFilterNotifier(cap, []string{"secret/prod"})
	msg := baseMsg("secret/prod/db")
	if err := f.Send(msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cap.msgs) != 1 {
		t.Fatalf("expected 1 forwarded message, got %d", len(cap.msgs))
	}
}

func TestFilterNotifier_NoMatchSuppresses(t *testing.T) {
	cap := &captureNotifier{}
	f, _ := NewFilterNotifier(cap, []string{"secret/prod"})
	if err := f.Send(baseMsg("secret/staging/db")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cap.msgs) != 0 {
		t.Fatal("expected message to be suppressed")
	}
}

func TestFilterNotifier_MultiplePrefix(t *testing.T) {
	cap := &captureNotifier{}
	f, _ := NewFilterNotifier(cap, []string{"secret/prod", "secret/staging"})
	f.Send(baseMsg("secret/prod/db"))
	f.Send(baseMsg("secret/staging/api"))
	f.Send(baseMsg("secret/dev/svc"))
	if len(cap.msgs) != 2 {
		t.Fatalf("expected 2 forwarded messages, got %d", len(cap.msgs))
	}
}

func TestFilterNotifier_PropagatesError(t *testing.T) {
	cap := &captureNotifier{err: errors.New("send failed")}
	f, _ := NewFilterNotifier(cap, []string{"secret/"})
	err := f.Send(baseMsg("secret/prod"))
	if err == nil {
		t.Fatal("expected error to propagate")
	}
}
