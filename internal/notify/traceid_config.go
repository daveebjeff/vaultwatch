package notify

import "fmt"

// TraceIDConfig holds configuration for constructing a TraceIDNotifier via
// the builder helpers used in cmd/vaultwatch/main.go.
type TraceIDConfig struct {
	// Header is the label key written to each message. Defaults to "trace_id".
	Header string `yaml:"header"`

	// PropagateEnvVar, when non-empty, names an environment variable whose
	// value is used as the trace ID for every Send call (useful in CI/CD
	// pipelines that set a build ID).
	PropagateEnvVar string `yaml:"propagate_env_var"`
}

// Validate returns an error if the config contains invalid values.
func (c TraceIDConfig) Validate() error {
	if len(c.Header) > 64 {
		return fmt.Errorf("traceid: header label key must be 64 characters or fewer, got %d", len(c.Header))
	}
	return nil
}

// Apply builds a TraceIDNotifier wrapping inner using the receiver config.
func (c TraceIDConfig) Apply(inner Notifier) (*TraceIDNotifier, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return NewTraceIDNotifier(inner, c.Header)
}
