package notify

import (
	"errors"
	"strings"
	"testing"
	"time"
)

type captureNotifier struct {
	got Message
	err error
}

func (c *captureNotifier) Send(m Message) error {
	c.got = m
	return c.err
}

func TestNewTemplateNotifier_NilInner(t *testing.T) {
	_, err := NewTemplateNotifier(nil, "")
	if err == nil {
		t.Fatal("expected error for nil inner")
	}
}

func TestNewTemplateNotifier_BadTemplate(t *testing.T) {
	_, err := NewTemplateNotifier(&captureNotifier{}, "{{.Unclosed")
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestTemplateNotifier_DefaultTemplate(t *testing.T) {
	cap := &captureNotifier{}
	n, err := NewTemplateNotifier(cap, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	msg := Message{
		Path:      "secret/db",
		Status:    StatusExpired,
		ExpiresAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
	}
	if err := n.Send(msg); err != nil {
		t.Fatalf("Send error: %v", err)
	}
	if !strings.Contains(cap.got.Path, "secret/db") {
		t.Errorf("rendered path missing secret path: %q", cap.got.Path)
	}
	if !strings.Contains(cap.got.Path, "EXPIRED") {
		t.Errorf("rendered path missing status: %q", cap.got.Path)
	}
}

func TestTemplateNotifier_CustomTemplate(t *testing.T) {
	cap := &captureNotifier{}
	n, err := NewTemplateNotifier(cap, "alert:{{.Path}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	msg := Message{Path: "secret/api", Status: StatusExpiringSoon}
	if err := n.Send(msg); err != nil {
		t.Fatalf("Send error: %v", err)
	}
	if cap.got.Path != "alert:secret/api" {
		t.Errorf("unexpected rendered path: %q", cap.got.Path)
	}
}

func TestTemplateNotifier_InnerError(t *testing.T) {
	cap := &captureNotifier{err: errors.New("downstream fail")}
	n, _ := NewTemplateNotifier(cap, "")
	err := n.Send(Message{Path: "x", ExpiresAt: time.Now()})
	if err == nil {
		t.Fatal("expected error from inner notifier")
	}
}

func TestTemplateNotifier_OriginalMessagePreserved(t *testing.T) {
	cap := &captureNotifier{}
	n, _ := NewTemplateNotifier(cap, "rendered")
	orig := Message{Path: "secret/orig", Status: StatusExpired, ExpiresAt: time.Now()}
	_ = n.Send(orig)
	if cap.got.Status != StatusExpired {
		t.Errorf("status not preserved: %v", cap.got.Status)
	}
}

func TestTemplateNotifier_ExpiresAtPreserved(t *testing.T) {
	cap := &captureNotifier{}
	n, _ := NewTemplateNotifier(cap, "rendered")
	expected := time.Date(2025, 1, 15, 9, 30, 0, 0, time.UTC)
	orig := Message{Path: "secret/token", Status: StatusExpiringSoon, ExpiresAt: expected}
	_ = n.Send(orig)
	if !cap.got.ExpiresAt.Equal(expected) {
		t.Errorf("ExpiresAt not preserved: got %v, want %v", cap.got.ExpiresAt, expected)
	}
}
