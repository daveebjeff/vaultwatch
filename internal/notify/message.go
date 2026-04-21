package notify

import "time"

// Status represents the expiration state of a Vault secret or lease.
type Status string

const (
	// StatusExpired indicates the secret has already expired.
	StatusExpired Status = "expired"
	// StatusExpiringSoon indicates the secret will expire within the warn window.
	StatusExpiringSoon Status = "expiring_soon"
	// StatusOK indicates the secret is healthy and not approaching expiration.
	StatusOK Status = "ok"
)

// Message carries the information sent to every Notifier implementation.
type Message struct {
	// Path is the Vault secret or lease path that triggered the alert.
	Path string
	// Status is the current expiration state of the secret.
	Status Status
	// ExpiresAt is the absolute time at which the secret or lease expires.
	ExpiresAt time.Time
	// Labels holds arbitrary key-value metadata attached to the message.
	// Middleware notifiers (e.g. LabelNotifier, MetadataNotifier) may add
	// entries here before forwarding the message downstream.
	Labels map[string]string
}
