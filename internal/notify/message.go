package notify

import "time"

// Status represents the expiration state of a secret.
type Status string

const (
	// StatusOK means the secret is not expiring soon.
	StatusOK Status = "ok"
	// StatusExpiringSoon means the secret will expire within the warn window.
	StatusExpiringSoon Status = "expiring_soon"
	// StatusExpired means the secret has already expired.
	StatusExpired Status = "expired"
)

// Message carries the alert details sent to notifiers.
type Message struct {
	// Summary is a human-readable one-line description.
	Summary string
	// SecretPath is the Vault path of the secret.
	SecretPath string
	// Status is the expiration status.
	Status Status
	// ExpiresAt is the time the secret expires.
	ExpiresAt time.Time
	// Details holds optional extra key/value metadata.
	Details map[string]string
}
