// Package notify — NormalizeNotifier
//
// NormalizeNotifier is a middleware Notifier that sanitizes the Body field
// of every outbound Message before forwarding it to the wrapped inner Notifier.
//
// Normalization steps (applied in order):
//  1. Trim leading and trailing whitespace (including newlines and tabs).
//  2. Collapse every internal run of whitespace characters into a single space.
//  3. Optionally convert the body to lower-case (opt-in via WithLowerCase).
//
// Usage:
//
//	n, err := notify.NewNormalizeNotifier(inner)
//	// or with lower-case conversion:
//	n, err := notify.NewNormalizeNotifier(inner, notify.WithLowerCase())
//
// NormalizeNotifier is useful when downstream systems are sensitive to
// inconsistent spacing or casing, such as log parsers or dedup keys.
package notify
