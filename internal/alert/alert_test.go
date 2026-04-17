package alert

import (
	"testing"
	"time"
)

func TestEvaluate_Expired(t *testing.T) {
	expired := time.Now().Add(-1 * time.Minute)
	a := Evaluate("secret/db", expired, 24*time.Hour)
	if a == nil {
		t.Fatal("expected alert for expired secret, got nil")
	}
	if a.Level != LevelCritical {
		t.Errorf("expected CRITICAL, got %s", a.Level)
	}
}

func TestEvaluate_ExpiringSoon(t *testing.T) {
	soon := time.Now().Add(1 * time.Hour)
	a := Evaluate("secret/db", soon, 24*time.Hour)
	if a == nil {
		t.Fatal("expected alert for expiring-soon secret, got nil")
	}
	if a.Level != LevelWarning {
		t.Errorf("expected WARNING, got %s", a.Level)
	}
}

func TestEvaluate_NotExpiring(t *testing.T) {
	future := time.Now().Add(72 * time.Hour)
	a := Evaluate("secret/db", future, 24*time.Hour)
	if a != nil {
		t.Errorf("expected no alert, got %+v", a)
	}
}

func TestEvaluate_ZeroTime(t *testing.T) {
	a := Evaluate("secret/db", time.Time{}, 24*time.Hour)
	if a != nil {
		t.Errorf("expected no alert for zero time, got %+v", a)
	}
}

func TestLogNotifier_Send(t *testing.T) {
	n := &LogNotifier{}
	a := Alert{
		Level:      LevelWarning,
		SecretPath: "secret/test",
		ExpiresAt:  time.Now().Add(2 * time.Hour),
		Message:    "expires soon",
	}
	if err := n.Send(a); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
