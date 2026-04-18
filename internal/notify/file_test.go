package notify

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestNewFileNotifier_EmptyPath(t *testing.T) {
	_, err := NewFileNotifier("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNewFileNotifier_Valid(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "vaultwatch-*.log")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	n, err := NewFileNotifier(tmp.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.path != tmp.Name() {
		t.Errorf("path mismatch")
	}
}

func TestFileNotifier_Send_WritesJSON(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/alerts.log"

	n, err := NewFileNotifier(path)
	if err != nil {
		t.Fatal(err)
	}

	msg := Message{
		Status:     StatusExpiringSoon,
		SecretPath: "secret/db/password",
		ExpiresAt:  time.Now().Add(10 * time.Minute),
	}
	if err := n.Send(msg); err != nil {
		t.Fatalf("Send error: %v", err)
	}

	f, _ := os.Open(path)
	defer f.Close()
	var got Message
	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		t.Fatal("no line written")
	}
	if err := json.Unmarshal(scanner.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.SecretPath != msg.SecretPath {
		t.Errorf("expected %q, got %q", msg.SecretPath, got.SecretPath)
	}
	if got.Status != msg.Status {
		t.Errorf("expected status %q, got %q", msg.Status, got.Status)
	}
}

func TestFileNotifier_Send_BadPath(t *testing.T) {
	n := &FileNotifier{path: "/nonexistent-dir/alerts.log"}
	msg := Message{Status: StatusExpired, SecretPath: "x", ExpiresAt: time.Now()}
	if err := n.Send(msg); err == nil {
		t.Fatal("expected error for bad path")
	}
}
