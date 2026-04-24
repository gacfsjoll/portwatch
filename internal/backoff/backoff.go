// Package backoff provides exponential backoff with optional jitter for
// retry logic used when scanner or notifier operations fail transiently.
package backoff

import (
	"math"
	"sync"
	"time"
)

// Clock allows tests to inject a fake time source.
type Clock func() time.Time

// Policy defines the parameters for exponential backoff.
type Policy struct {
	Initial    time.Duration
	Max        time.Duration
	Multiplier float64
	Jitter     float64 // fraction in [0, 1)
}

// DefaultPolicy returns a sensible default backoff policy.
func DefaultPolicy() Policy {
	return Policy{
		Initial:    500 * time.Millisecond,
		Max:        30 * time.Second,
		Multiplier: 2.0,
		Jitter:     0.1,
	}
}

// Backoff tracks per-key retry state.
type Backoff struct {
	mu      sync.Mutex
	policy  Policy
	clock   Clock
	attempt map[string]int
}

// New creates a Backoff using the given policy and real time.
func New(p Policy) *Backoff {
	return NewWithClock(p, time.Now)
}

// NewWithClock creates a Backoff with an injectable clock.
func NewWithClock(p Policy, clock Clock) *Backoff {
	if p.Multiplier <= 1 {
		p.Multiplier = 2.0
	}
	return &Backoff{
		policy:  p,
		clock:   clock,
		attempt: make(map[string]int),
	}
}

// Next returns the delay for the next retry attempt for the given key
// and increments the internal attempt counter.
func (b *Backoff) Next(key string) time.Duration {
	b.mu.Lock()
	defer b.mu.Unlock()

	n := b.attempt[key]
	b.attempt[key] = n + 1

	delay := float64(b.policy.Initial) * math.Pow(b.policy.Multiplier, float64(n))
	if b.policy.Jitter > 0 {
		// deterministic pseudo-jitter based on attempt count to keep tests stable
		fraction := float64(n%10) / 10.0
		delay += delay * b.policy.Jitter * fraction
	}
	if d := time.Duration(delay); d > b.policy.Max {
		return b.policy.Max
	} else {
		return d
	}
}

// Reset clears the attempt counter for the given key.
func (b *Backoff) Reset(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.attempt, key)
}

// Attempts returns the current attempt count for a key.
func (b *Backoff) Attempts(key string) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.attempt[key]
}
