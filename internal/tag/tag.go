// Package tag provides a lightweight key-value tagging system for port events.
// Tags can be attached to alerts to carry contextual metadata such as environment,
// host role, or custom annotations through the notification pipeline.
package tag

import (
	"fmt"
	"sort"
	"strings"
)

// Set is an immutable collection of string tags in "key=value" form.
type Set struct {
	tags map[string]string
}

// New constructs a Set from a slice of "key=value" pairs.
// Duplicate keys are last-write-wins. An error is returned if any
// pair is malformed.
func New(pairs []string) (Set, error) {
	m := make(map[string]string, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			return Set{}, fmt.Errorf("tag: invalid pair %q: must be key=value", p)
		}
		m[parts[0]] = parts[1]
	}
	return Set{tags: m}, nil
}

// Get returns the value for key and whether it was present.
func (s Set) Get(key string) (string, bool) {
	v, ok := s.tags[key]
	return v, ok
}

// All returns a sorted copy of all key=value pairs.
func (s Set) All() []string {
	out := make([]string, 0, len(s.tags))
	for k, v := range s.tags {
		out = append(out, k+"="+v)
	}
	sort.Strings(out)
	return out
}

// Len returns the number of tags in the set.
func (s Set) Len() int { return len(s.tags) }

// String returns a comma-separated representation of all tags.
func (s Set) String() string {
	return strings.Join(s.All(), ",")
}

// Merge returns a new Set combining s and other. Keys in other override s.
func (s Set) Merge(other Set) Set {
	m := make(map[string]string, len(s.tags)+len(other.tags))
	for k, v := range s.tags {
		m[k] = v
	}
	for k, v := range other.tags {
		m[k] = v
	}
	return Set{tags: m}
}
