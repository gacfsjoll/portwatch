// Package redact provides utilities for scrubbing sensitive values from
// alert payloads and log lines before they are written to external sinks.
//
// A Redactor holds a set of literal strings and regular-expression patterns
// that should be replaced with a fixed placeholder so that secrets such as
// API keys or bearer tokens never appear in outbound notifications.
package redact

import (
	"regexp"
	"strings"
	"sync"
)

const placeholder = "[REDACTED]"

// Redactor scrubs sensitive strings from text.
type Redactor struct {
	mu       sync.RWMutex
	literals []string
	patterns []*regexp.Regexp
}

// New returns a Redactor pre-loaded with the supplied literal strings.
func New(literals ...string) *Redactor {
	r := &Redactor{}
	for _, l := range literals {
		if l != "" {
			r.literals = append(r.literals, l)
		}
	}
	return r
}

// AddPattern compiles expr and registers it as an additional redaction rule.
// The entire match is replaced by the placeholder.
func (r *Redactor) AddPattern(expr string) error {
	re, err := regexp.Compile(expr)
	if err != nil {
		return err
	}
	r.mu.Lock()
	r.patterns = append(r.patterns, re)
	r.mu.Unlock()
	return nil
}

// AddLiteral registers a literal string that should be redacted.
func (r *Redactor) AddLiteral(s string) {
	if s == "" {
		return
	}
	r.mu.Lock()
	r.literals = append(r.literals, s)
	r.mu.Unlock()
}

// Scrub returns a copy of text with all registered literals and pattern
// matches replaced by the placeholder.
func (r *Redactor) Scrub(text string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, l := range r.literals {
		text = strings.ReplaceAll(text, l, placeholder)
	}
	for _, re := range r.patterns {
		text = re.ReplaceAllString(text, placeholder)
	}
	return text
}
