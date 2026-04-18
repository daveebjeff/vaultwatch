package notify

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func exampleMsg() Message {
	return Message{
		Level:     LevelWarning,
		Secret:    "secret/db/password",
		ExpiresAt: time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC),
		Details:   "expires in 24h",
	}
}

func TestLogNotifier_Send_ContainsFields(t *testing.T) {
	var buf bytes.Buffer
	ln := &LogNotifier{Out: &buf}
	if err := ln.Send(exampleMsg()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"WARNING", "secret/db/password", "expires in 24h", "2025-06-01"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q; got: %s", want, out)
		}
	}
}

func TestLogNotifier_Send_WriterError(t *testing.T) {
	ln := &LogNotifier{Out: &errorWriter{}}
	if err := ln.Send(exampleMsg()); err == nil {
		t.Fatal("expected error from failing writer")
	}
}

// errorWriter always returns an error.
type errorWriter struct{}

func (e *errorWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("write error")
}
