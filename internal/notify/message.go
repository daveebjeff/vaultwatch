package notify

import "time"

// Status represents the alert status for a secret.
type Status string

const (
	// StatusExpired indicates the secret has already expired.
	StatusExpired Status = "EXPIRED"
	// StatusExpiringSoon indicates the secret will expire soon.
	StatusExpiringSoon Status = "EXPIRING_SOON"
	// StatusOK indicates the secret is healthy.
	StatusOK Status = "OK"
)

// Message holds the data passed to every Notifier.
type Message struct {
	// SecretPath is the Vault path of the secret.
	SecretPath string
	// Status is the current alert status.
	Status Status
	// Expiry is the time at which the secret expires.
	Expiry time.Time
	// Detail is a human-readable description of the alert.
	Detail string
}

// Notifier is the interface implemented by all notification backends.
type Notifier interface {
	Send(msg Message) error
}
