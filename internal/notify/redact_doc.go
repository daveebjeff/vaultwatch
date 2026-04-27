// Package notify — RedactNotifier
//
// RedactNotifier scrubs sensitive values from alert message bodies before
// they are forwarded to the downstream notifier. It applies one or more
// compiled regular expressions and replaces matches with a configurable
// replacement string (default: "[REDACTED]").
//
// # Usage
//
//	patterns, err := notify.CompilePatterns([]string{
//		`(?i)password=\S+`,
//		`s\.[A-Za-z0-9]{24,}`,
//	})
//	if err != nil { ... }
//
//	n, err := notify.NewRedactNotifier(inner, patterns, "[REDACTED]")
//
// # Defaults
//
// NewDefaultRedactNotifier applies a set of built-in patterns that cover
// common Vault token formats and generic key=value secret fields.
//
// # Notes
//
//   - Only the Body field of the Message is modified; Labels and other
//     metadata are left unchanged.
//   - Patterns are applied in the order they are provided.
package notify
