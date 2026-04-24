// Package retention provides a policy for pruning old history and audit
// entries so that on-disk state does not grow without bound.
package retention

import (
	"time"
)

// Policy describes how long different classes of records should be kept.
type Policy struct {
	// MaxAge is the maximum age of a record before it is eligible for pruning.
	MaxAge time.Duration
	// MaxEntries is the maximum number of records to retain (0 = unlimited).
	MaxEntries int
}

// DefaultPolicy returns a sensible out-of-the-box retention policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAge:     30 * 24 * time.Hour, // 30 days
		MaxEntries: 10_000,
	}
}

// Pruner applies a Policy to a slice of time-stamped records and returns
// only those that should be kept.
type Pruner struct {
	policy Policy
	now    func() time.Time
}

// New creates a Pruner using the supplied policy and the real wall clock.
func New(p Policy) *Pruner {
	return NewWithClock(p, time.Now)
}

// NewWithClock creates a Pruner with an injectable clock (useful in tests).
func NewWithClock(p Policy, now func() time.Time) *Pruner {
	return &Pruner{policy: p, now: now}
}

// Apply filters entries older than MaxAge and, if MaxEntries > 0, trims the
// result to the most-recent MaxEntries records.
//
// entries must be ordered oldest-first; the returned slice preserves that
// ordering.
func (pr *Pruner) Apply(entries []time.Time) []time.Time {
	cutoff := pr.now().Add(-pr.policy.MaxAge)

	filtered := entries[:0:0]
	for _, t := range entries {
		if !t.Before(cutoff) {
			filtered = append(filtered, t)
		}
	}

	if pr.policy.MaxEntries > 0 && len(filtered) > pr.policy.MaxEntries {
		filtered = filtered[len(filtered)-pr.policy.MaxEntries:]
	}

	return filtered
}

// ShouldPrune reports whether the given timestamp is older than the policy
// MaxAge and therefore eligible for removal.
func (pr *Pruner) ShouldPrune(t time.Time) bool {
	cutoff := pr.now().Add(-pr.policy.MaxAge)
	return t.Before(cutoff)
}
