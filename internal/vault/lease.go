package vault

import (
	"fmt"
	"time"
)

// LeaseInfo holds metadata about a secret lease.
type LeaseInfo struct {
	Path      string
	LeaseID   string
	Duration  time.Duration
	Renewable bool
	ExpiresAt time.Time
}

// GetLeaseInfo extracts lease information from a Vault secret.
func GetLeaseInfo(path string, secret interface {
	GetLeaseDuration() int
	GetLeaseID() string
	IsRenewable() bool
}) LeaseInfo {
	duration := time.Duration(secret.GetLeaseDuration()) * time.Second
	return LeaseInfo{
		Path:      path,
		LeaseID:   secret.GetLeaseID(),
		Duration:  duration,
		Renewable: secret.IsRenewable(),
		ExpiresAt: time.Now().Add(duration),
	}
}

// IsExpiringSoon returns true if the lease expires within the given threshold.
func (l LeaseInfo) IsExpiringSoon(threshold time.Duration) bool {
	return time.Until(l.ExpiresAt) <= threshold
}

// String returns a human-readable summary of the lease.
func (l LeaseInfo) String() string {
	return fmt.Sprintf("path=%s lease_id=%s expires_at=%s renewable=%v",
		l.Path, l.LeaseID, l.ExpiresAt.Format(time.RFC3339), l.Renewable)
}
