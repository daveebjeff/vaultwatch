package notify

import "time"

// CooldownConfig holds configuration for building a CooldownNotifier via
// dependency injection or config-driven wiring.
type CooldownConfig struct {
	// Cooldown is the quiet period enforced between successive forwards for
	// the same secret path. Must be greater than zero.
	Cooldown time.Duration `yaml:"cooldown" json:"cooldown"`
}

// Validate returns an error if the configuration is invalid.
func (c CooldownConfig) Validate() error {
	if c.Cooldown <= 0 {
		return ErrZeroCooldown
	}
	return nil
}

// Build constructs a CooldownNotifier wrapping inner using the receiver's
// configuration. Validate is called before construction.
func (c CooldownConfig) Build(inner Notifier) (*CooldownNotifier, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return NewCooldownNotifier(inner, c.Cooldown)
}
