package notify

import "fmt"

// TemplateConfig holds configuration for a TemplateNotifier loaded from
// vaultwatch config files.
type TemplateConfig struct {
	// Template is the Go text/template string used to render alerts.
	// Leave empty to use DefaultTemplate.
	Template string `yaml:"template" json:"template"`
}

// Validate checks that the template string, if provided, is parseable.
func (c TemplateConfig) Validate() error {
	if c.Template == "" {
		return nil
	}
	_, err := newParsedTemplate(c.Template)
	if err != nil {
		return fmt.Errorf("template config: %w", err)
	}
	return nil
}

// Build wraps inner with a TemplateNotifier using the configured template.
func (c TemplateConfig) Build(inner Notifier) (*TemplateNotifier, error) {
	return NewTemplateNotifier(inner, c.Template)
}

// newParsedTemplate is a thin helper so Validate and NewTemplateNotifier
// share the same parse path without importing text/template twice at the
// call site.
func newParsedTemplate(tmplStr string) (interface{}, error) {
	import_text_template_via_build_constraint_avoidance := func() error {
		_, err := NewTemplateNotifier(&multiNotifier{}, tmplStr)
		return err
	}
	return nil, import_text_template_via_build_constraint_avoidance()
}
