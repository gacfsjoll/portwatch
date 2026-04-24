// Package deadman implements a dead-man's switch that fires an alert when
// the monitor has not completed a successful scan within a configurable
// deadline. This guards against silent failures where the daemon is running
// but scans are silently stalling.
package deadman

import (
	"sync"
	"time"
)

// Notifier is the callback invoked when the dead-man's switch trips.
type Notifier func(lastSeen time.Time, elapsed time.Duration)

// Switch is a dead-man's switch that trips if Reset is not called within
// the configured deadline.
type Switch struct {
	mu       sync.Mutex
	deadline time.Duration
	lastSeen time.Time
	notify   Notifier
	clock    func() time.Time
	stop     chan struct{}
	wg       sync.WaitGroup
}

// New creates a new Switch that will call notify if Reset is not called
// within deadline. The switch begins monitoring immediately.
func New(deadline time.Duration, notify Notifier) *Switch {
	return NewWithClock(deadline, notify, time.Now)
}

// NewWithClock creates a Switch with an injectable clock (useful for testing).
func NewWithClock(deadline time.Duration, notify Notifier, clock func() time.Time) *Switch {
	s := &Switch{
		deadline: deadline,
		notify:   notify,
		clock:    clock,
		lastSeen: clock(),
		stop:     make(chan struct{}),
	}
	s.wg.Add(1)
	go s.run()
	return s
}

// Reset records a successful heartbeat, preventing the switch from tripping.
func (s *Switch) Reset() {
	s.mu.Lock()
	s.lastSeen = s.clock()
	s.mu.Unlock()
}

// Stop shuts down the background goroutine.
func (s *Switch) Stop() {
	close(s.stop)
	s.wg.Wait()
}

// LastSeen returns the time of the most recent Reset call.
func (s *Switch) LastSeen() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lastSeen
}

func (s *Switch) run() {
	defer s.wg.Done()
	ticker := time.NewTicker(s.deadline / 2)
	defer ticker.Stop()
	for {
		select {
		case <-s.stop:
			return
		case <-ticker.C:
			s.mu.Lock()
			last := s.lastSeen
			s.mu.Unlock()
			now := s.clock()
			if now.Sub(last) >= s.deadline {
				s.notify(last, now.Sub(last))
			}
		}
	}
}
