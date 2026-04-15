package monitor

import (
	"log"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Change represents a detected port state change.
type Change struct {
	Port   int
	Old    scanner.PortState
	New    scanner.PortState
}

// Monitor periodically scans ports and reports changes.
type Monitor struct {
	scanner  *scanner.Scanner
	interval time.Duration
	previous map[int]scanner.PortState
	OnChange func(Change)
}

// New creates a Monitor for the given port range and poll interval.
func New(startPort, endPort int, interval time.Duration) *Monitor {
	return &Monitor{
		scanner:  scanner.NewScanner(startPort, endPort),
		interval: interval,
		previous: make(map[int]scanner.PortState),
		OnChange: defaultOnChange,
	}
}

// Start begins the monitoring loop. It blocks until stop is closed.
func (m *Monitor) Start(stop <-chan struct{}) error {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	// Run an initial scan so we have a baseline.
	if err := m.scan(); err != nil {
		return err
	}

	for {
		select {
		case <-ticker.C:
			if err := m.scan(); err != nil {
				log.Printf("portwatch: scan error: %v", err)
			}
		case <-stop:
			return nil
		}
	}
}

func (m *Monitor) scan() error {
	states, err := m.scanner.Scan()
	if err != nil {
		return err
	}

	current := make(map[int]scanner.PortState, len(states))
	for _, s := range states {
		current[s.Port] = s
	}

	for port, newState := range current {
		if oldState, seen := m.previous[port]; seen {
			if oldState.Open != newState.Open {
				m.OnChange(Change{Port: port, Old: oldState, New: newState})
			}
		}
	}

	m.previous = current
	return nil
}

func defaultOnChange(c Change) {
	if c.New.Open {
		log.Printf("[ALERT] port %d opened", c.Port)
	} else {
		log.Printf("[ALERT] port %d closed", c.Port)
	}
}
