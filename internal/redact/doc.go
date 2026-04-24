// Package redact provides a thread-safe Redactor that removes sensitive
// values — such as API tokens, passwords, or bearer credentials — from
// strings before they are forwarded to external alert sinks or written to
// log files.
//
// # Usage
//
//	r := redact.New("my-api-key", "db-password")
//	_ = r.AddPattern(`Bearer\s+\S+`)
//	clean := r.Scrub(payload)
//
// Both literal replacements and compiled regular-expression patterns are
// applied in registration order.  Every match is replaced by the fixed
// placeholder string "[REDACTED]".
package redact
