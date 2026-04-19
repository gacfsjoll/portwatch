// Package debounce prevents alert storms by suppressing repeated notifications
// for the same port event within a configurable window.
package debounce

import (
	"sync"
	"time"
)

// Clock allows time to be injected for testing.
type Clock func() time.Time

// Debouncer tracks the last alert time per port and suppresses duplicates
// within the cooldown window.
type Debouncer struct {
	mu       sync.Mutex
	window   time.Duration
	clock    Clock
	lastSeen map[string]time.Time
}

// New creates a Debouncer with the given window duration.
func New(window time.Duration) *Debouncer {
	return NewWithClock(window, time.Now)
}

// NewWithClock creates a Debouncer with an injectable clock.
func NewWithClock(window time.Duration, clock Clock) *Debouncer {
	return &Debouncer{
		window:   window,
		clock:    clock,
		lastSeen: make(map[string]time.Time),
	}
}

// Allow returns true if the event for the given key should be forwarded.
// It returns false if an identical event was already seen within the window.
func (d *Debouncer) Allow(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.clock()
	if last, ok := d.lastSeen[key]; ok {
		if now.Sub(last) < d.window {
			return false
		}
	}
	d.lastSeen[key] = now
	return true
}

// Reset removes the debounce record for the given key.
func (d *Debouncer) Reset(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.lastSeen, key)
}

// Flush clears all debounce state.
func (d *Debouncer) Flush() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.lastSeen = make(map[string]time.Time)
}
