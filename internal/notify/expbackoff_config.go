package notify

import "time"

// ExpBackoffConfig holds YAML-serialisable configuration for
// ExpBackoffNotifier.
type ExpBackoffConfig struct {
	// InitDelay is the initial retry delay (e.g. "200ms").
	InitDelay time.Duration `yaml:"init_delay"`
	// MaxDelay caps the per-attempt delay (e.g. "30s").
	MaxDelay time.Duration `yaml:"max_delay"`
	// Attempts is the maximum number of total send attempts.
	Attempts int `yaml:"attempts"`
}

// DefaultExpBackoffConfig returns a sensible default configuration.
func DefaultExpBackoffConfig() ExpBackoffConfig {
	return ExpBackoffConfig{
		InitDelay: 500 * time.Millisecond,
		MaxDelay:  30 * time.Second,
		Attempts:  4,
	}
}

// BuildExpBackoffNotifier constructs an ExpBackoffNotifier from cfg,
// wrapping inner.
func BuildExpBackoffNotifier(inner Notifier, cfg ExpBackoffConfig) (*ExpBackoffNotifier, error) {
	if cfg.InitDelay <= 0 {
		cfg.InitDelay = DefaultExpBackoffConfig().InitDelay
	}
	if cfg.MaxDelay <= 0 {
		cfg.MaxDelay = DefaultExpBackoffConfig().MaxDelay
	}
	if cfg.Attempts < 1 {
		cfg.Attempts = DefaultExpBackoffConfig().Attempts
	}
	return NewExpBackoffNotifier(inner, cfg.InitDelay, cfg.MaxDelay, cfg.Attempts)
}
