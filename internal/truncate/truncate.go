// Package truncate provides utilities for capping alert message payloads
// to a configurable maximum byte length, ensuring downstream notifiers
// (webhooks, log sinks) are not overwhelmed by unexpectedly large messages.
package truncate

import "fmt"

const (
	// DefaultMaxBytes is the default payload size cap (4 KiB).
	DefaultMaxBytes = 4096

	suffix = "... [truncated]"
)

// Truncator caps string payloads at a maximum byte length.
type Truncator struct {
	maxBytes int
}

// New returns a Truncator with the given maximum byte length.
// If maxBytes is <= 0 the DefaultMaxBytes value is used.
func New(maxBytes int) *Truncator {
	if maxBytes <= 0 {
		maxBytes = DefaultMaxBytes
	}
	return &Truncator{maxBytes: maxBytes}
}

// Apply returns s unchanged when len(s) <= maxBytes, otherwise it returns
// a prefix of s (in bytes) followed by the truncation suffix so that the
// total length does not exceed maxBytes.
func (t *Truncator) Apply(s string) string {
	if len(s) <= t.maxBytes {
		return s
	}

	avail := t.maxBytes - len(suffix)
	if avail <= 0 {
		// maxBytes is smaller than the suffix itself; just cut hard.
		return s[:t.maxBytes]
	}

	// Trim to a valid UTF-8 boundary.
	cut := avail
	for cut > 0 && !isRuneStart(s[cut]) {
		cut--
	}

	return s[:cut] + suffix
}

// Applyf is a convenience wrapper that formats according to a format
// specifier and then applies the byte-length cap.
func (t *Truncator) Applyf(format string, args ...any) string {
	return t.Apply(fmt.Sprintf(format, args...))
}

// MaxBytes returns the configured cap.
func (t *Truncator) MaxBytes() int { return t.maxBytes }

// isRuneStart reports whether b is the first byte of a UTF-8 sequence.
func isRuneStart(b byte) bool {
	// Continuation bytes have the form 10xxxxxx.
	return b&0xC0 != 0x80
}
