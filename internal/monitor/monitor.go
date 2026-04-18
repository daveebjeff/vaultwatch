package monitor

import (
	"context"
	"log"
	"time"

	"github.com/example/vaultwatch/internal/alert"
	"github.com/example/vaultwatch/internal/vault"
)

// SecretPath represents a secret path to monitor.
type SecretPath struct {
	Path     string
	LeaseTTL time.Duration
}

// Monitor polls Vault for secret lease info and triggers alerts.
type Monitor struct {
	client   *vault.Client
	notifier alert.Notifier
	paths    []SecretPath
	interval time.Duration
	warnBefore time.Duration
}

// New creates a new Monitor.
func New(client *vault.Client, notifier alert.Notifier, paths []SecretPath, interval, warnBefore time.Duration) *Monitor {
	return &Monitor{
		client:     client,
		notifier:   notifier,
		paths:      paths,
		interval:   interval,
		warnBefore: warnBefore,
	}
}

// Run starts the monitoring loop until ctx is cancelled.
func (m *Monitor) Run(ctx context.Context) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()
	log.Printf("monitor: starting, interval=%s warn_before=%s", m.interval, m.warnBefore)
	for {
		select {
		case <-ctx.Done():
			log.Println("monitor: stopping")
			return
		case <-ticker.C:
			m.check()
		}
	}
}

func (m *Monitor) check() {
	for _, sp := range m.paths {
		info, err := vault.GetLeaseInfo(m.client, sp.Path)
		if err != nil {
			log.Printf("monitor: error fetching lease for %s: %v", sp.Path, err)
			continue
		}
		events := alert.Evaluate(info, m.warnBefore)
		for _, ev := range events {
			if err := m.notifier.Send(ev); err != nil {
				log.Printf("monitor: failed to send alert for %s: %v", sp.Path, err)
			}
		}
	}
}
