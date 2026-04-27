package notify

import "fmt"

// RedactConfig holds the configuration needed to build a RedactNotifier
// from external sources such as a YAML config file.
type RedactConfig struct {
	// Patterns is a list of raw regular expression strings.
	Patterns []string `yaml:"patterns" json:"patterns"`
	// Replacement is the string used to replace matched text.
	// Defaults to "[REDACTED]" when empty.
	Replacement string `yaml:"replacement" json:"replacement"`
	// UseDefaults, when true, prepends the built-in secret patterns to
	// any patterns provided in Patterns.
	UseDefaults bool `yaml:"use_defaults" json:"use_defaults"`
}

// Build constructs a RedactNotifier from the config, wrapping inner.
func (c RedactConfig) Build(inner Notifier) (*RedactNotifier, error) {
	var patterns []*regexp.Regexp

	if c.UseDefaults {
		patterns = append(patterns, defaultRedactPatterns...)
	}

	if len(c.Patterns) > 0 {
		compiled, err := CompilePatterns(c.Patterns)
		if err != nil {
			return nil, fmt.Errorf("redact config: %w", err)
		}
		patterns = append(patterns, compiled...)
	}

	if len(patterns) == 0 {
		// Fall back to defaults when nothing is specified.
		patterns = defaultRedactPatterns
	}

	return NewRedactNotifier(inner, patterns, c.Replacement)
}
