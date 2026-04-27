// Package limiter provides a token-bucket style alert limiter that caps the
// number of notifications emitted per port within a rolling time window.
package limiter

import (
	"sync"
	"time"
)

// Clock allows tests to inject a fake time source.
type Clock func() time.Time

// entry tracks the token count and window start for a single port.
type entry struct {
	count     int
	windowStart time.Time
}

// Limiter caps alert emissions per port to MaxAlerts within Window.
type Limiter struct {
	MaxAlerts int
	Window    time.Duration
	clock     Clock
	mu        sync.Mutex
	buckets   map[int]*entry
}

// New creates a Limiter with the given cap and rolling window.
func New(maxAlerts int, window time.Duration) *Limiter {
	return NewWithClock(maxAlerts, window, time.Now)
}

// NewWithClock creates a Limiter with a custom clock (useful in tests).
func NewWithClock(maxAlerts int, window time.Duration, clock Clock) *Limiter {
	if maxAlerts < 1 {
		maxAlerts = 1
	}
	return &Limiter{
		MaxAlerts: maxAlerts,
		Window:    window,
		clock:     clock,
		buckets:   make(map[int]*entry),
	}
}

// Allow returns true if the alert for the given port should be forwarded.
// It increments the counter for the port and resets the window when expired.
func (l *Limiter) Allow(port int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	e, ok := l.buckets[port]
	if !ok || now.Sub(e.windowStart) >= l.Window {
		l.buckets[port] = &entry{count: 1, windowStart: now}
		return true
	}
	if e.count >= l.MaxAlerts {
		return false
	}
	e.count++
	return true
}

// Reset clears the token bucket for a specific port.
func (l *Limiter) Reset(port int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, port)
}

// Stats returns the current hit count and window-start time for a port.
// Returns zero values if the port has no active bucket.
func (l *Limiter) Stats(port int) (count int, since time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if e, ok := l.buckets[port]; ok {
		return e.count, e.windowStart
	}
	return 0, time.Time{}
}
