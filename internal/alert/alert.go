package alert

import (
	"fmt"
	"log"
	"time"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo    Level = "INFO"
	LevelWarning Level = "WARNING"
	LevelCritical Level = "CRITICAL"
)

// Alert represents a secret expiration alert.
type Alert struct {
	Level     Level
	SecretPath string
	ExpiresAt  time.Time
	Message    string
}

// Notifier defines the interface for sending alerts.
type Notifier interface {
	Send(alert Alert) error
}

// LogNotifier sends alerts to stdout via the standard logger.
type LogNotifier struct{}

// Send logs the alert to stdout.
func (l *LogNotifier) Send(a Alert) error {
	log.Printf("[%s] secret=%s expires_at=%s message=%s",
		a.Level, a.SecretPath, a.ExpiresAt.Format(time.RFC3339), a.Message)
	return nil
}

// Evaluate inspects a secret's expiry and returns an Alert if action is needed.
// warnThreshold is how far in advance to warn before expiry.
func Evaluate(secretPath string, expiresAt time.Time, warnThreshold time.Duration) *Alert {
	now := time.Now()

	if expiresAt.Is	return nil
	}

	ttl := expiresAt.Sub(now)

	switch {
	case ttl <= 0:
		return &tLevel:      LevelCritical,
			SecretPath: secretPath,
			ExpiresAt:  expiresAt,
			Message:    "secret hast}
	case ttl <= warnThreshold:
		return &Alert{
			Level:      LevelWarning,
			SecretPath: secretPath,
			ExpiresAt:  expiresAt,
			Message:    fmt.Sprintf("secret expires in %s", ttl.Round(time.Second)),
		}
	default:
		return nil
	}
}
