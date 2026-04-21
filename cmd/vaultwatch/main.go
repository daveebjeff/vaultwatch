// Package main is the entry point for the vaultwatch CLI tool.
// It wires together configuration, Vault client, monitor, and notifiers
// to provide continuous secret expiration monitoring and alerting.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/youorg/vaultwatch/internal/alert"
	"github.com/youorg/vaultwatch/internal/config"
	"github.com/youorg/vaultwatch/internal/monitor"
	"github.com/youorg/vaultwatch/internal/notify"
	"github.com/youorg/vaultwatch/internal/vault"
)

con = "0.1.0"

func main() {
	var (
		configPath  = flag.String("config", "config.yaml", "path to configuration file")
		showVersion = flag.Bool("version", false, "print version and exit")
		dryRun      = flag.Bool("dry-run", false, "evaluate secrets once and exit without looping")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("vaultwatch v%s\n", version)
		os.Exit(0)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	vaultClient, err := vault.NewClient(cfg.Vault.Address, cfg.Vault.Token)
	if err != nil {
		log.Fatalf("failed to create vault client: %v", err)
	}

	// Build the notifier chain: log notifier as the default sink.
	var notifier notify.Notifier = notify.NewLogNotifier(os.Stdout)

	// Wrap with retry for transient failures.
	notifier = notify.NewRetryNotifier(notifier, 3, 2*time.Second)

	// Wrap with deduplication to suppress repeated identical alerts.
	notifier = notify.NewDedupNotifier(notifier)

	// Wrap with rate limiting to avoid alert storms (1 alert per path per minute).
	notifier = notify.NewRateLimitNotifier(notifier, time.Minute)

	evaluator := alert.Evaluate

	monCfg := monitor.Config{
		Interval:   cfg.Monitor.Interval,
		WarnBefore: cfg.Monitor.WarnBefore,
		Paths:      toSecretPaths(cfg.Monitor.Paths),
	}
	if err := monCfg.Validate(); err != nil {
		log.Fatalf("invalid monitor config: %v", err)
	}

	mon := monitor.New(monCfg, vaultClient, evaluator, notifier)

	if *dryRun {
		log.Println("dry-run mode: running single evaluation pass")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := mon.RunOnce(ctx); err != nil {
			log.Fatalf("dry-run evaluation failed: %v", err)
		}
		log.Println("dry-run complete")
		os.Exit(0)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Printf("vaultwatch v%s starting (interval=%s, warn_before=%s, paths=%d)",
		version, monCfg.Interval, monCfg.WarnBefore, len(monCfg.Paths))

	if err := mon.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalf("monitor exited with error: %v", err)
	}

	log.Println("vaultwatch stopped")
}

// toSecretPaths converts config path entries into monitor.SecretPath values.
func toSecretPaths(entries []config.PathEntry) []monitor.SecretPath {
	paths := make([]monitor.SecretPath, 0, len(entries))
	for _, e := range entries {
		paths = append(paths, monitor.SecretPath{
			Path:  e.Path,
			Alias: e.Alias,
		})
	}
	return paths
}
