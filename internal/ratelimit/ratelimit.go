// Package ratelimit provides alert rate-limiting to suppress duplicate
// notifications for the same port within a configurable cooldown window.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter tracks the last alert time per port and suppresses duplicates
// within the cooldown period.
type Limiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[uint16]time.Time
	now      func() time.Time
}

// New creates a Limiter with the given cooldown duration.
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		cooldown: cooldown,
		last:     make(map[uint16]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if an alert for port should be sent, and records the
// current time as the last-alert time. Returns false if the port is still
// within the cooldown window.
func (l *Limiter) Allow(port uint16) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	if t, ok := l.last[port]; ok {
		if now.Sub(t) < l.cooldown {
			return false
		}
	}
	l.last[port] = now
	return true
}

// Reset clears the rate-limit record for a specific port.
func (l *Limiter) Reset(port uint16) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, port)
}

// ResetAll clears all recorded alert times.
func (l *Limiter) ResetAll() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.last = make(map[uint16]time.Time)
}
