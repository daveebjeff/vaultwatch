package vault

import (
	"testing"
	"time"
)

type mockSecret struct {
	duration  int
	leaseID   string
	renewable bool
}

func (m mockSecret) GetLeaseDuration() int { return m.duration }
func (m mockSecret) GetLeaseID() string    { return m.leaseID }
func (m mockSecret) IsRenewable() bool     { return m.renewable }

func TestGetLeaseInfo(t *testing.T) {
	secret := mockSecret{duration: 3600, leaseID: "lease-abc", renewable: true}
	info := GetLeaseInfo("secret/myapp", secret)

	if info.Path != "secret/myapp" {
		t.Errorf("expected path secret/myapp, got %s", info.Path)
	}
	if info.LeaseID != "lease-abc" {
		t.Errorf("expected lease-abc, got %s", info.LeaseID)
	}
	if info.Duration != 3600*time.Second {
		t.Errorf("unexpected duration: %v", info.Duration)
	}
	if !info.Renewable {
		t.Error("expected renewable to be true")
	}
}

func TestIsExpiringSoon(t *testing.T) {
	now := time.Now()

	expiring := LeaseInfo{ExpiresAt: now.Add(5 * time.Minute)}
	if !expiring.IsExpiringSoon(10 * time.Minute) {
		t.Error("expected lease to be expiring soon")
	}

	fresh := LeaseInfo{ExpiresAt: now.Add(2 * time.Hour)}
	if fresh.IsExpiringSoon(10 * time.Minute) {
		t.Error("expected lease to not be expiring soon")
	}
}

func TestLeaseInfoString(t *testing.T) {
	l := LeaseInfo{
		Path:      "secret/db",
		LeaseID:   "id-123",
		Renewable: false,
		ExpiresAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	s := l.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
