package notify

import (
	"fmt"
	"os/exec"
)

// ScriptNotifier runs an external script/command when an alert fires,
// passing status and secret path as arguments.
type ScriptNotifier struct {
	scriptPath string
}

// NewScriptNotifier creates a ScriptNotifier for the given executable path.
func NewScriptNotifier(scriptPath string) (*ScriptNotifier, error) {
	if scriptPath == "" {
		return nil, fmt.Errorf("script: path must not be empty")
	}
	return &ScriptNotifier{scriptPath: scriptPath}, nil
}

// Send executes the script with status, secret path, and expiry as arguments.
func (s *ScriptNotifier) Send(msg Message) error {
	expiry := msg.ExpiresAt.UTC().Format("2006-01-02T15:04:05Z")
	cmd := exec.Command(s.scriptPath, string(msg.Status), msg.SecretPath, expiry) //nolint:gosec
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("script: execution failed: %w (output: %s)", err, string(output))
	}
	return nil
}
