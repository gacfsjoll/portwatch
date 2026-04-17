// Package suppress provides a suppression list for port change alerts.
// Ports added to the suppression list will not trigger notifications
// for a configurable duration.
package suppress

import (
	"sync"
	"time"
)

// Entry holds suppression metadata for a single port.
type Entry struct {
	Port      int
	Until     time.Time
	Reason    string
}

// List manages a set of suppressed ports.
type List struct {
	mu      sync.RWMutex
	entries map[int]Entry
	now     func() time.Time
}

// New returns an initialised suppression List.
func New() *List {
	return &List{
		entries: make(map[int]Entry),
		now:     time.Now,
	}
}

// Suppress silences alerts for port for the given duration.
func (l *List) Suppress(port int, d time.Duration, reason string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries[port] = Entry{
		Port:   port,
		Until:  l.now().Add(d),
		Reason: reason,
	}
}

// IsSuppressed reports whether port is currently suppressed.
func (l *List) IsSuppressed(port int) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	e, ok := l.entries[port]
	if !ok {
		return false
	}
	return l.now().Before(e.Until)
}

// Remove lifts suppression for port immediately.
func (l *List) Remove(port int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.entries, port)
}

// Expire removes all entries whose suppression window has passed.
func (l *List) Expire() {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	for p, e := range l.entries {
		if !now.Before(e.Until) {
			delete(l.entries, p)
		}
	}
}

// All returns a snapshot of all active suppression entries.
func (l *List) All() []Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Entry, 0, len(l.entries))
	for _, e := range l.entries {
		out = append(out, e)
	}
	return out
}
