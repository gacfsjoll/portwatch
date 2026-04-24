// Package rollup batches multiple port-change events within a short window
// and emits a single summarised alert, reducing noise when many ports change
// at the same time (e.g. during a service restart).
package rollup

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Flusher is called with the accumulated events when the window closes.
type Flusher func(events []alert.Event)

// Roller accumulates events and flushes them after a quiet window.
type Roller struct {
	mu      sync.Mutex
	window  time.Duration
	events  []alert.Event
	timer   *time.Timer
	flush   Flusher
	nowFunc func() time.Time
}

// New returns a Roller that waits window after the last Add call before
// invoking flush with all accumulated events.
func New(window time.Duration, flush Flusher) *Roller {
	return &Roller{
		window:  window,
		flush:   flush,
		nowFunc: time.Now,
	}
}

// Add appends an event to the current batch and (re)starts the flush timer.
func (r *Roller) Add(e alert.Event) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.events = append(r.events, e)

	if r.timer != nil {
		r.timer.Stop()
	}
	r.timer = time.AfterFunc(r.window, r.drain)
}

// Flush forces an immediate flush regardless of the window, useful on
// shutdown to ensure no events are lost.
func (r *Roller) Flush() {
	r.mu.Lock()
	if r.timer != nil {
		r.timer.Stop()
		r.timer = nil
	}
	r.mu.Unlock()
	r.drain()
}

// drain is called by the timer or Flush; it emits accumulated events.
func (r *Roller) drain() {
	r.mu.Lock()
	if len(r.events) == 0 {
		r.mu.Unlock()
		return
	}
	batch := make([]alert.Event, len(r.events))
	copy(batch, r.events)
	r.events = r.events[:0]
	r.mu.Unlock()

	r.flush(batch)
}
