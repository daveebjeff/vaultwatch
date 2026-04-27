package notify

import (
	"fmt"
	"time"
)

// WatermarkConfig holds the configuration for building a WatermarkNotifier.
type WatermarkConfig struct {
	// Threshold is the remaining TTL below which the notification fires.
	// Must be a positive duration string, e.g. "24h", "6h30m".
	Threshold string `yaml:"threshold" json:"threshold"`
}

// Validate checks that the WatermarkConfig is well-formed.
func (c WatermarkConfig) Validate() error {
	if c.Threshold == "" {
		return fmt.Errorf("watermark config: threshold must not be empty")
	}
	d, err := time.ParseDuration(c.Threshold)
	if err != nil {
		return fmt.Errorf("watermark config: invalid threshold %q: %w", c.Threshold, err)
	}
	if d <= 0 {
		return fmt.Errorf("watermark config: threshold must be positive, got %s", c.Threshold)
	}
	return nil
}

// Build constructs a WatermarkNotifier from the config and the provided inner notifier.
func (c WatermarkConfig) Build(inner Notifier) (*WatermarkNotifier, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	d, _ := time.ParseDuration(c.Threshold) // already validated
	return NewWatermarkNotifier(inner, d)
}
