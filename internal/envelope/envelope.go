// Package envelope wraps alert events with metadata such as hostname,
// process ID, and a monotonically increasing sequence number before they
// are handed off to notifiers. This ensures every outbound alert carries
// enough context to be correlated without consulting additional sources.
package envelope

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

// Envelope wraps an arbitrary payload with delivery metadata.
type Envelope struct {
	// Seq is a process-scoped monotonic counter that increments for every
	// envelope produced during the lifetime of the daemon.
	Seq uint64 `json:"seq"`

	// Hostname is the machine that produced this envelope.
	Hostname string `json:"hostname"`

	// PID is the OS process ID of the portwatch daemon.
	PID int `json:"pid"`

	// At is the UTC timestamp at which the envelope was created.
	At time.Time `json:"at"`

	// Payload is the wrapped value (e.g. an alert.Event).
	Payload any `json:"payload"`
}

// String returns a compact human-readable representation.
func (e Envelope) String() string {
	return fmt.Sprintf("[#%d %s pid=%d %s]", e.Seq, e.Hostname, e.PID, e.At.Format(time.RFC3339))
}

// Wrapper produces Envelopes. It caches hostname resolution and owns the
// sequence counter so callers never need to manage either.
type Wrapper struct {
	hostname string
	pid      int
	counter  atomic.Uint64
	now      func() time.Time
}

// New returns a Wrapper ready for use. Hostname resolution failures are
// silently replaced with "unknown" so the daemon never crashes on startup.
func New() *Wrapper {
	h, err := os.Hostname()
	if err != nil {
		h = "unknown"
	}
	return &Wrapper{
		hostname: h,
		pid:      os.Getpid(),
		now:      time.Now,
	}
}

// newWithClock is used in tests to inject a deterministic clock.
func newWithClock(now func() time.Time) *Wrapper {
	w := New()
	w.now = now
	return w
}

// Wrap increments the sequence counter and returns a new Envelope containing
// payload.
func (w *Wrapper) Wrap(payload any) Envelope {
	return Envelope{
		Seq:      w.counter.Add(1),
		Hostname: w.hostname,
		PID:      w.pid,
		At:       w.now().UTC(),
		Payload:  payload,
	}
}
