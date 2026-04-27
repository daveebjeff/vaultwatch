package notify

import (
	"fmt"
	"time"
)

// EscalationConfig holds configuration for building an EscalationNotifier
// from a YAML/env-driven setup.
type EscalationConfig struct {
	// Timeout is how long an unacknowledged alert waits before escalation.
	// Parsed from a duration string, e.g. "15m".
	Timeout string `yaml:"timeout" mapstructure:"timeout"`
}

// Validate returns an error if the config is unusable.
func (c EscalationConfig) Validate() error {
	if c.Timeout == "" {
		return fmt.Errorf("escalation_config: timeout must not be empty")
	}
	if _, err := time.ParseDuration(c.Timeout); err != nil {
		return fmt.Errorf("escalation_config: invalid timeout %q: %w", c.Timeout, err)
	}
	return nil
}

// ParsedTimeout returns the timeout as a time.Duration.
func (c EscalationConfig) ParsedTimeout() (time.Duration, error) {
	return time.ParseDuration(c.Timeout)
}

// Build constructs an EscalationNotifier from the config and provided notifiers.
// Both primary and secondary must be non-nil.
func (c EscalationConfig) Build(primary, secondary Notifier) (*EscalationNotifier, error) {
	if primary == nil {
		return nil, fmt.Errorf("escalation_config: primary notifier must not be nil")
	}
	if secondary == nil {
		return nil, fmt.Errorf("escalation_config: secondary notifier must not be nil")
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	d, err := c.ParsedTimeout()
	if err != nil {
		return nil, err
	}
	return NewEscalationNotifier(primary, secondary, d)
}
