package notify

// NormalizeConfig holds configuration options for the NormalizeNotifier.
type NormalizeConfig struct {
	// LowerCase controls whether the message body is converted to lower case.
	// Defaults to false.
	LowerCase bool `yaml:"lower_case" json:"lower_case"`

	// TrimSpace controls whether leading and trailing whitespace is stripped
	// from the message body. Defaults to true.
	TrimSpace bool `yaml:"trim_space" json:"trim_space"`

	// CollapseSpaces controls whether runs of internal whitespace are
	// collapsed to a single space. Defaults to true.
	CollapseSpaces bool `yaml:"collapse_spaces" json:"collapse_spaces"`
}

// DefaultNormalizeConfig returns a NormalizeConfig with sensible defaults.
func DefaultNormalizeConfig() NormalizeConfig {
	return NormalizeConfig{
		LowerCase:      false,
		TrimSpace:      true,
		CollapseSpaces: true,
	}
}

// options converts the config into a slice of functional options suitable
// for passing to NewNormalizeNotifier.
func (c NormalizeConfig) options() []func(*NormalizeNotifier) {
	var opts []func(*NormalizeNotifier)
	if c.LowerCase {
		opts = append(opts, WithLowerCase())
	}
	return opts
}

// BuildNormalizeNotifier constructs a NormalizeNotifier from a config and an
// inner Notifier. It applies all options derived from the config.
func BuildNormalizeNotifier(inner Notifier, cfg NormalizeConfig) (*NormalizeNotifier, error) {
	return NewNormalizeNotifier(inner, cfg.options()...)
}
