package monitor_test

import (
	"testing"
	"time"

	"github.com/example/vaultwatch/internal/monitor"
)

func TestValidate_Valid(t *testing.T) {
	cfg := &monitor.Config{
		Interval:   30 * time.Second,
		WarnBefore: 24 * time.Hour,
		Paths:      []monitor.SecretPath{{Path: "secret/db"}},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidate_ZeroInterval(t *testing.T) {
	cfg := &monitor.Config{
		Interval:   0,
		WarnBefore: 24 * time.Hour,
		Paths:      []monitor.SecretPath{{Path: "secret/db"}},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestValidate_NoPaths(t *testing.T) {
	cfg := &monitor.Config{
		Interval:   time.Minute,
		WarnBefore: time.Hour,
		Paths:      []monitor.SecretPath{},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestValidate_EmptyPathEntry(t *testing.T) {
	cfg := &monitor.Config{
		Interval:   time.Minute,
		WarnBefore: time.Hour,
		Paths:      []monitor.SecretPath{{Path: ""}},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty path entry")
	}
}

func TestValidate_ZeroWarnBefore(t *testing.T) {
	cfg := &monitor.Config{
		Interval:   time.Minute,
		WarnBefore: 0,
		Paths:      []monitor.SecretPath{{Path: "secret/app"}},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero warn_before")
	}
}

func TestValidate_NegativeInterval(t *testing.T) {
	cfg := &monitor.Config{
		Interval:   -1 * time.Second,
		WarnBefore: time.Hour,
		Paths:      []monitor.SecretPath{{Path: "secret/app"}},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative interval")
	}
}
