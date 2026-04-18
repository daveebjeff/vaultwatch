package notify

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// FileNotifier appends JSON-encoded alert messages to a file on disk.
type FileNotifier struct {
	mu   sync.Mutex
	path string
}

// NewFileNotifier creates a FileNotifier that writes to the given file path.
// The file is created if it does not exist.
func NewFileNotifier(path string) (*FileNotifier, error) {
	if path == "" {
		return nil, fmt.Errorf("file notifier: path must not be empty")
	}
	// Probe that we can open/create the file.
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("file notifier: open %s: %w", path, err)
	}
	f.Close()
	return &FileNotifier{path: path}, nil
}

// Send appends a JSON line representing msg to the configured file.
func (f *FileNotifier) Send(msg Message) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	file, err := os.OpenFile(f.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("file notifier: open: %w", err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	if err := enc.Encode(msg); err != nil {
		return fmt.Errorf("file notifier: encode: %w", err)
	}
	return nil
}
