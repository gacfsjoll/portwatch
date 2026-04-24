// Package cooldown provides per-port alert suppression based on a
// minimum quiet period between repeated notifications for the same port.
//
// Unlike ratelimit (which enforces a fixed cooldown after any alert) and
// throttle (which caps total alerts in a window), cooldown focuses on
// preventing alert storms caused by flapping ports — ports that open and
// close repeatedly within a short time span.
package cooldown

import (
	"sync"
	"time"
)

// Clock is a function that returns the current time. Swappable in tests.
type Clock func() time.Time

// Tracker records the last alert time per port and enforces a quiet period.
type Tracker struct {
	mu      sync.Mutex
	last    map[int]time.Time
	period  time.Duration
	clock   Clock
}

// New returns a Tracker that suppresses repeated alerts for the same port
// within the given quiet period.
func New(period time.Duration) *Tracker {
	return NewWithClock(period, time.Now)
}

// NewWithClock returns a Tracker with a custom clock, useful in tests.
func NewWithClock(period time.Duration, clock Clock) *Tracker {
	return &Tracker{
		last:   make(map[int]time.Time),
		period: period,
		clock:  clock,
	}
}

// Allow returns true if enough time has elapsed since the last alert for
// the given port, and records the current time as the new last-alert time.
func (t *Tracker) Allow(port int) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.clock()
	if last, ok := t.last[port]; ok {
		if now.Sub(last) < t.period {
			return false
		}
	}
	t.last[port] = now
	return true
}

// Reset clears the recorded alert time for the given port, allowing the
// next alert through immediately.
func (t *Tracker) Reset(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, port)
}

// Active returns the number of ports currently within their quiet period.
func (t *Tracker) Active() int {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.clock()
	count := 0
	for _, last := range t.last {
		if now.Sub(last) < t.period {
			count++
		}
	}
	return count
}
