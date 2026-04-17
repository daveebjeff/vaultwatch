package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "vaultwatch-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTempConfig(t, `
vault:
  address: "https://vault.example.com"
  token: "s.testtoken"
  interval: 30s
alerting:
  warn_threshold: 48h
  critical_threshold: 12h
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Address != "https://vault.example.com" {
		t.Errorf("expected address, got %q", cfg.Vault.Address)
	}
	if cfg.Vault.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", cfg.Vault.Interval)
	}
	if cfg.Alerting.CriticalThreshold != 12*time.Hour {
		t.Errorf("expected 12h critical threshold, got %v", cfg.Alerting.CriticalThreshold)
	}
}

func TestLoad_Defaults(t *testing.T) {
	path := writeTempConfig(t, `
vault:
  address: "https://vault.example.com"
  token: "s.testtoken"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Interval != 60*time.Second {
		t.Errorf("expected default 60s interval, got %v", cfg.Vault.Interval)
	}
	if cfg.Alerting.WarnThreshold != 7*24*time.Hour {
		t.Errorf("expected default 7d warn threshold, got %v", cfg.Alerting.WarnThreshold)
	}
}

func TestLoad_MissingAddress(t *testing.T) {
	path := writeTempConfig(t, `
vault:
  token: "s.testtoken"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestLoad_TokenFromEnv(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "s.envtoken")
	path := writeTempConfig(t, `
vault:
  address: "https://vault.example.com"
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Token != "s.envtoken" {
		t.Errorf("expected token from env, got %q", cfg.Vault.Token)
	}
}
