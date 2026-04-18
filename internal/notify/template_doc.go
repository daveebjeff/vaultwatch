// Package notify provides notification backends for vaultwatch.
//
// TemplateNotifier
//
// TemplateNotifier wraps any Notifier and renders each Message through a
// Go text/template before forwarding it. The rendered output is stored in
// Message.Path so downstream notifiers (e.g. LogNotifier, FileNotifier)
// receive a human-readable string without needing their own formatting logic.
//
// Example:
//
//	 n, _ := notify.NewTemplateNotifier(inner,
//	     "[{{.Status}}] {{.Path}} — expires {{.ExpiresAt.Format \"2006-01-02\"}}",
//	 )
//	 n.Send(msg)
package notify
