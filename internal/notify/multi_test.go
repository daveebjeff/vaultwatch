package notify

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

type fakeNotifier struct {
	called bool
	err    error
}

func (f *fakeNotifier) Send(_ Message) error {
	f.called = true
	return f.err
}

func TestMultiNotifier_AllCalled(t *testing.T) {
	a, b := &fakeNotifier{}, &fakeNotifier{}
	mn := NewMultiNotifier(a, b)
	if err := mn.Send(exampleMsg()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !a.called || !b.called {
		t.Error("expected both notifiers to be called")
	}
}

func TestMultiNotifier_CollectsErrors(t *testing.T) {
	a := &fakeNotifier{err: errors.New("fail a")}
	b := &fakeNotifier{err: errors.New("fail b")}
	mn := NewMultiNotifier(a, b)
	err := mn.Send(exampleMsg())
	if err == nil {
		t.Fatal("expected combined error")
	}
	if !strings.Contains(err.Error(), "2 backend(s) failed") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestMultiNotifier_Add(t *testing.T) {
	var buf bytes.Buffer
	mn := NewMultiNotifier()
	mn.Add(&LogNotifier{Out: &buf})
	if err := mn.Send(exampleMsg()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected output from added notifier")
	}
}
