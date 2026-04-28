package notify

import "fmt"

// EnrichConfig holds configuration for building an EnrichNotifier via
// BuildEnrichNotifier. All fields are optional; omitting a field leaves
// the corresponding enrichment behaviour at its default.
type EnrichConfig struct {
	// AddHostname, when true, stamps the local hostname into the message
	// labels under the key "hostname".
	AddHostname bool `yaml:"add_hostname" json:"add_hostname"`

	// AddSeverity, when true, derives a severity string from the message
	// status and stamps it into the message labels under the key "severity".
	AddSeverity bool `yaml:"add_severity" json:"add_severity"`

	// AddEnvironment stamps the given string into the message labels under
	// the key "environment". Ignored when empty.
	AddEnvironment string `yaml:"environment" json:"environment"`

	// ExtraLabels are arbitrary key/value pairs that are stamped into every
	// outgoing message. Values must not be empty strings.
	ExtraLabels map[string]string `yaml:"extra_labels" json:"extra_labels"`
}

// Validate returns an error if the configuration contains invalid values.
func (c EnrichConfig) Validate() error {
	for k, v := range c.ExtraLabels {
		if k == "" {
			return fmt.Errorf("enrich_config: extra_labels contains an empty key")
		}
		if v == "" {
			return fmt.Errorf("enrich_config: extra_labels[%q] has an empty value", k)
		}
	}
	return nil
}

// BuildEnrichNotifier wraps inner with an EnrichNotifier configured
// according to cfg. If no enrichment options are enabled the original
// inner notifier is returned unwrapped so that no overhead is added.
func BuildEnrichNotifier(inner Notifier, cfg EnrichConfig) (Notifier, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	opts := []EnrichOption{}

	if cfg.AddHostname {
		opts = append(opts, WithHostname())
	}

	if cfg.AddSeverity {
		opts = append(opts, WithSeverityLabel())
	}

	if cfg.AddEnvironment != "" {
		opts = append(opts, WithEnvironment(cfg.AddEnvironment))
	}

	for k, v := range cfg.ExtraLabels {
		opts = append(opts, WithStaticLabel(k, v))
	}

	if len(opts) == 0 {
		return inner, nil
	}

	return NewEnrichNotifier(inner, opts...)
}
