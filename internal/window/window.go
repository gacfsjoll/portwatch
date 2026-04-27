// Package window provides a sliding time-window counter for tracking
// event frequency over a rolling duration. It is safe for concurrent use.
package window

import (
	"sync"
	"time"
)

// Clock allows injecting a fake time source in tests.
type Clock func() time.Time

// Counter tracks how many events have occurred within a sliding window.
type Counter struct {
	mu       sync.Mutex
	window   time.Duration
	clock    Clock
	buckets  []entry
}

type entry struct {
	at    time.Time
	count int
}

// New returns a Counter that tracks events within the given sliding window
// duration using the real wall clock.
func New(window time.Duration) *Counter {
	return NewWithClock(window, time.Now)
}

// NewWithClock returns a Counter using the supplied clock function.
func NewWithClock(window time.Duration, clock Clock) *Counter {
	return &Counter{
		window:  window,
		clock:   clock,
		buckets: make([]entry, 0, 64),
	}
}

// Add records n events at the current time.
func (c *Counter) Add(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.clock()
	c.buckets = append(c.buckets, entry{at: now, count: n})
	c.evict(now)
}

// Total returns the sum of all events recorded within the current window.
func (c *Counter) Total() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evict(c.clock())
	total := 0
	for _, b := range c.buckets {
		total += b.count
	}
	return total
}

// Reset discards all recorded events.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.buckets = c.buckets[:0]
}

// evict removes entries that have fallen outside the window. Must be called
// with c.mu held.
func (c *Counter) evict(now time.Time) {
	cutoff := now.Add(-c.window)
	i := 0
	for i < len(c.buckets) && c.buckets[i].at.Before(cutoff) {
		i++
	}
	c.buckets = c.buckets[i:]
}
