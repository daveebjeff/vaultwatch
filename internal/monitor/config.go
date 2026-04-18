package monitor

import (
	"fmt"
	"time"
)

// Config holds monitor-specific configuration.
type Config struct {
	// Interval is how often to poll Vault.
	Interval time.Duration
	// WarnBefore is the time window before expiry to trigger a warning.
	WarnBefore time.Duration
	// Paths is the list of secret paths to watch.
	Paths []SecretPath
}

// Validate checks that Config values are sensible.
func (c *Config) Validate() error {
	if c.Interval <= 0 {
		return fmt.Errorf("monitor: interval must be positive, got %s", c.Interval)
	}
	if c.WarnBefore <= 0 {
		return fmt.Errorf("monitor: warn_before must be positive, got %s", c.WarnBefore)
	}
	if len(c.Paths) == 0 {
		return fmt.Errorf("monitor: at least one secret path must be specified")
	}
	for i, p := range c.Paths {
		if p.Path == "" {
			return fmt.Errorf("monitor: path at index %d is empty", i)
		}
	}
	return nil
}
