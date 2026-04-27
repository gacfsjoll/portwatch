// Package dedup provides event deduplication based on a rolling time window.
// Duplicate events for the same port and change type are suppressed until the
// deduplication window expires, preventing alert storms when a port flaps.
package dedup

import (
	"sync"
	"time"
)

// key uniquely identifies an alert event by port and direction.
type key struct {
	port      int
	direction string
}

// entry records when a key was last seen.
type entry struct {
	seenAt time.Time
}

// Deduplicator suppresses repeated events for the same port within a window.
type Deduplicator struct {
	mu     sync.Mutex
	window time.Duration
	seen   map[key]entry
	now    func() time.Time
}

// New returns a Deduplicator with the given deduplication window.
func New(window time.Duration) *Deduplicator {
	return NewWithClock(window, time.Now)
}

// NewWithClock returns a Deduplicator using a custom clock, useful for testing.
func NewWithClock(window time.Duration, now func() time.Time) *Deduplicator {
	return &Deduplicator{
		window: window,
		seen:   make(map[key]entry),
		now:    now,
	}
}

// Allow returns true if the event for the given port and direction should be
// forwarded. Returns false if an identical event was already seen within the
// deduplication window.
func (d *Deduplicator) Allow(port int, direction string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	k := key{port: port, direction: direction}

	if e, ok := d.seen[k]; ok {
		if now.Sub(e.seenAt) < d.window {
			return false
		}
	}

	d.seen[k] = entry{seenAt: now}
	return true
}

// Reset removes the deduplication record for the given port and direction,
// allowing the next event to pass through immediately.
func (d *Deduplicator) Reset(port int, direction string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.seen, key{port: port, direction: direction})
}

// Purge evicts all entries whose window has expired. It is safe to call
// periodically to reclaim memory in long-running daemons.
func (d *Deduplicator) Purge() {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	for k, e := range d.seen {
		if now.Sub(e.seenAt) >= d.window {
			delete(d.seen, k)
		}
	}
}
