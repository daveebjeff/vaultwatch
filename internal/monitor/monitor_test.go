package monitor_test

import (
	"context"
	"testing"
	"time"

	"github.com/example/vaultwatch/internal/monitor"
)

func TestNew_DefaultsSet(t *testing.T) {
	m := monitor.New(nil, nil, nil, 30*time.Second, 24*time.Hour)
	if m == nil {
		t.Fatal("expected non-nil monitor")
	}
}

func TestRun_CancelImmediately(t *testing.T) {
	m := monitor.New(nil, nil, []monitor.SecretPath{}, time.Hour, 24*time.Hour)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	done := make(chan struct{})
	go func() {
		m.Run(ctx)
		close(done)
	}()
	select {
	case <-done:
		// ok
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not stop after context cancellation")
	}
}

func TestSecretPath_Fields(t *testing.T) {
	sp := monitor.SecretPath{
		Path:     "secret/myapp/db",
		LeaseTTL: 72 * time.Hour,
	}
	if sp.Path != "secret/myapp/db" {
		t.Errorf("unexpected path: %s", sp.Path)
	}
	if sp.LeaseTTL != 72*time.Hour {
		t.Errorf("unexpected TTL: %s", sp.LeaseTTL)
	}
}
