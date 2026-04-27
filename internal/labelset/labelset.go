// Package labelset provides a key-value label store that can be attached
// to port events for richer alerting context (e.g. environment, owner, tier).
package labelset

import (
	"fmt"
	"sort"
	"strings"
)

// LabelSet is an immutable collection of key=value labels.
type LabelSet struct {
	labels map[string]string
}

// New creates a LabelSet from the provided key=value pairs.
// Pairs that do not contain exactly one '=' are silently skipped.
func New(pairs ...string) (*LabelSet, error) {
	ls := &LabelSet{labels: make(map[string]string, len(pairs))}
	for _, p := range pairs {
		k, v, ok := strings.Cut(p, "=")
		if !ok {
			return nil, fmt.Errorf("labelset: invalid pair %q: must be key=value", p)
		}
		k = strings.TrimSpace(k)
		if k == "" {
			return nil, fmt.Errorf("labelset: empty key in pair %q", p)
		}
		ls.labels[k] = strings.TrimSpace(v)
	}
	return ls, nil
}

// Get returns the value for key and whether it was found.
func (ls *LabelSet) Get(key string) (string, bool) {
	v, ok := ls.labels[key]
	return v, ok
}

// All returns a copy of all labels as a map.
func (ls *LabelSet) All() map[string]string {
	out := make(map[string]string, len(ls.labels))
	for k, v := range ls.labels {
		out[k] = v
	}
	return out
}

// String returns a deterministic, comma-separated representation.
func (ls *LabelSet) String() string {
	keys := make([]string, 0, len(ls.labels))
	for k := range ls.labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+ls.labels[k])
	}
	return strings.Join(parts, ",")
}

// Merge returns a new LabelSet combining ls with other.
// Keys from other overwrite keys from ls on conflict.
func (ls *LabelSet) Merge(other *LabelSet) *LabelSet {
	out := &LabelSet{labels: make(map[string]string, len(ls.labels)+len(other.labels))}
	for k, v := range ls.labels {
		out.labels[k] = v
	}
	for k, v := range other.labels {
		out.labels[k] = v
	}
	return out
}
