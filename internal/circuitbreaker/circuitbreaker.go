// Package circuitbreaker implements a simple circuit breaker that stops
// forwarding alerts when a downstream notifier experiences repeated failures.
// Once the failure threshold is reached the circuit opens and calls are
// rejected until a configurable recovery window has elapsed.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// ErrCircuitOpen is returned when the circuit is open and calls are rejected.
var ErrCircuitOpen = errors.New("circuit breaker is open")

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed State = iota
	StateOpen
)

// String returns a human-readable representation of the state.
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	default:
		return "unknown"
	}
}

// Breaker tracks consecutive failures and opens the circuit when the threshold
// is exceeded. It uses a pluggable clock so tests can control time.
type Breaker struct {
	mu           sync.Mutex
	threshold    int
	recovery     time.Duration
	clock        func() time.Time
	failures     int
	openedAt     time.Time
	state        State
}

// New creates a Breaker that opens after threshold consecutive failures and
// attempts recovery after the given duration.
func New(threshold int, recovery time.Duration) *Breaker {
	return NewWithClock(threshold, recovery, time.Now)
}

// NewWithClock creates a Breaker with an injectable clock for testing.
func NewWithClock(threshold int, recovery time.Duration, clock func() time.Time) *Breaker {
	return &Breaker{
		threshold: threshold,
		recovery:  recovery,
		clock:     clock,
	}
}

// Allow returns nil if the call should proceed, or ErrCircuitOpen if the
// circuit is open and the recovery window has not yet elapsed.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.state == StateOpen {
		if b.clock().Sub(b.openedAt) >= b.recovery {
			// Recovery window elapsed — move back to closed (half-open probe).
			b.state = StateClosed
			b.failures = 0
		} else {
			return ErrCircuitOpen
		}
	}
	return nil
}

// RecordSuccess resets the failure counter.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure increments the failure counter and opens the circuit if the
// threshold has been reached.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.threshold && b.state == StateClosed {
		b.state = StateOpen
		b.openedAt = b.clock()
	}
}

// CurrentState returns the current state of the breaker.
func (b *Breaker) CurrentState() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
