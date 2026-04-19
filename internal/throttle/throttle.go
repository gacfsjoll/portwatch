// Package throttle limits the number of alerts fired within a rolling
// time window to prevent alert storms during rapid port churn.
package throttle

import (
	"sync"
	"time"
)

// Clock allows injecting a fake time source in tests.
type Clock func() time.Time

// Throttle tracks per-port alert counts within a rolling window.
type Throttle struct {
	mu      sync.Mutex
	window  time.Duration
	maxHits int
	clock   Clock
	hits    map[int][]time.Time
}

// New returns a Throttle that allows at most maxHits alerts per port
// within the given rolling window.
func New(window time.Duration, maxHits int) *Throttle {
	return NewWithClock(window, maxHits, time.Now)
}

// NewWithClock is like New but accepts an injectable clock.
func NewWithClock(window time.Duration, maxHits int, clock Clock) *Throttle {
	return &Throttle{
		window:  window,
		maxHits: maxHits,
		clock:   clock,
		hits:    make(map[int][]time.Time),
	}
}

// Allow returns true if an alert for port should be fired now.
// It records the attempt regardless of the outcome.
func (t *Throttle) Allow(port int) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.clock()
	cutoff := now.Add(-t.window)

	prev := t.hits[port]
	var recent []time.Time
	for _, ts := range prev {
		if ts.After(cutoff) {
			recent = append(recent, ts)
		}
	}

	if len(recent) >= t.maxHits {
		t.hits[port] = recent
		return false
	}

	t.hits[port] = append(recent, now)
	return true
}

// Reset clears the hit history for a specific port.
func (t *Throttle) Reset(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.hits, port)
}

// Stats returns the number of recorded hits within the current window for port.
func (t *Throttle) Stats(port int) int {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.clock()
	cutoff := now.Add(-t.window)
	count := 0
	for _, ts := range t.hits[port] {
		if ts.After(cutoff) {
			count++
		}
	}
	return count
}
