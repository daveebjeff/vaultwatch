package notify

import (
	"bytes"
	"fmt"
	"text/template"
)

// DefaultTemplate is the default message template.
const DefaultTemplate = `[{{.Status}}] {{.Path}} expires at {{.ExpiresAt.Format "2006-01-02 15:04:05 UTC"}}`

// TemplateNotifier renders a Message using a Go text/template before
// forwarding the formatted string to an inner Notifier via a synthetic
// Message whose Path carries the rendered text.
type TemplateNotifier struct {
	inner    Notifier
	tmpl     *template.Template
}

// NewTemplateNotifier creates a TemplateNotifier with the provided template
// string. Use DefaultTemplate as tmplStr for the built-in format.
func NewTemplateNotifier(inner Notifier, tmplStr string) (*TemplateNotifier, error) {
	if inner == nil {
		return nil, fmt.Errorf("template notifier: inner notifier must not be nil")
	}
	if tmplStr == "" {
		tmplStr = DefaultTemplate
	}
	t, err := template.New("notify").Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("template notifier: parse template: %w", err)
	}
	return &TemplateNotifier{inner: inner, tmpl: t}, nil
}

// Send renders msg with the configured template and forwards the result.
func (n *TemplateNotifier) Send(msg Message) error {
	var buf bytes.Buffer
	if err := n.tmpl.Execute(&buf, msg); err != nil {
		return fmt.Errorf("template notifier: render: %w", err)
	}
	rendered := msg
	rendered.Path = buf.String()
	return n.inner.Send(rendered)
}
