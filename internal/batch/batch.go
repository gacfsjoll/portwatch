// Package batch provides a size- and time-bounded event accumulator.
// Events are collected until either a maximum batch size is reached or
// a flush interval elapses, whichever comes first. The flusher is
// called with the accumulated slice and the buffer is reset.
package batch

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Flusher is called with a non-empty slice of events when the batch
// is ready to be dispatched.
type Flusher func(events []alert.Event)

// Accumulator collects alert.Event values and flushes them in batches.
type Accumulator struct {
	mu       sync.Mutex
	buf      []alert.Event
	maxSize  int
	interval time.Duration
	flusher  Flusher
	ticker   *time.Ticker
	stop     chan struct{}
	done     chan struct{}
}

// New creates an Accumulator that flushes when buf reaches maxSize events
// or when interval elapses, whichever comes first.
// maxSize must be >= 1; interval must be > 0.
func New(maxSize int, interval time.Duration, flusher Flusher) *Accumulator {
	if maxSize < 1 {
		maxSize = 1
	}
	if interval <= 0 {
		interval = time.Second
	}
	a := &Accumulator{
		buf:      make([]alert.Event, 0, maxSize),
		maxSize:  maxSize,
		interval: interval,
		flusher:  flusher,
		ticker:   time.NewTicker(interval),
		stop:     make(chan struct{}),
		done:     make(chan struct{}),
	}
	go a.run()
	return a
}

// Add appends an event to the current batch. If the batch reaches maxSize
// it is flushed immediately.
func (a *Accumulator) Add(e alert.Event) {
	a.mu.Lock()
	a.buf = append(a.buf, e)
	ready := len(a.buf) >= a.maxSize
	a.mu.Unlock()

	if ready {
		a.Flush()
	}
}

// Flush drains the current buffer and calls the flusher synchronously.
// It is a no-op when the buffer is empty.
func (a *Accumulator) Flush() {
	a.mu.Lock()
	if len(a.buf) == 0 {
		a.mu.Unlock()
		return
	}
	events := make([]alert.Event, len(a.buf))
	copy(events, a.buf)
	a.buf = a.buf[:0]
	a.mu.Unlock()

	a.flusher(events)
}

// Stop halts the background ticker and performs a final flush.
func (a *Accumulator) Stop() {
	a.ticker.Stop()
	close(a.stop)
	<-a.done
	a.Flush()
}

func (a *Accumulator) run() {
	defer close(a.done)
	for {
		select {
		case <-a.ticker.C:
			a.Flush()
		case <-a.stop:
			return
		}
	}
}
