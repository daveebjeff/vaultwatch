package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Vault   VaultConfig   `yaml:"vault"`
	Alerting AlertingConfig `yaml:"alerting"`
}

type VaultConfig struct {
	Address   string        `yaml:"address"`
	Token     string        `yaml:"token"`
	Namespace string        `yaml:"namespace"`
	Interval  time.Duration `yaml:"interval"`
}

type AlertingConfig struct {
	WarnThreshold     time.Duration `yaml:"warn_threshold"`
	CriticalThreshold time.Duration `yaml:"critical_threshold"`
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening config file: %w", err)
	}
	defer f.Close()

	cfg := &Config{}
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("decoding config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Vault.Address == "" {
		return fmt.Errorf("vault.address is required")
	}
	if c.Vault.Token == "" {
		c.Vault.Token = os.Getenv("VAULT_TOKEN")
	}
	if c.Vault.Token == "" {
		return fmt.Errorf("vault.token is required (or set VAULT_TOKEN env var)")
	}
	if c.Vault.Interval == 0 {
		c.Vault.Interval = 60 * time.Second
	}
	if c.Alerting.WarnThreshold == 0 {
		c.Alerting.WarnThreshold = 7 * 24 * time.Hour
	}
	if c.Alerting.CriticalThreshold == 0 {
		c.Alerting.CriticalThreshold = 24 * time.Hour
	}
	return nil
}
