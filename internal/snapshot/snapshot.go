// Package snapshot provides diffing between two sets of port states.
package snapshot

import "github.com/user/portwatch/internal/scanner"

// Diff represents the changes between two port snapshots.
type Diff struct {
	Opened []scanner.PortState
	Closed []scanner.PortState
}

// Compare returns the difference between a previous and current set of port states.
// Ports present in current but not previous are Opened; ports in previous but not current are Closed.
func Compare(previous, current []scanner.PortState) Diff {
	prev := index(previous)
	curr := index(current)

	var diff Diff

	for port, state := range curr {
		if _, existed := prev[port]; !existed {
			diff.Opened = append(diff.Opened, state)
		}
	}

	for port, state := range prev {
		if _, exists := curr[port]; !exists {
			diff.Closed = append(diff.Closed, state)
		}
	}

	return diff
}

// HasChanges reports whether the diff contains any opened or closed ports.
func (d Diff) HasChanges() bool {
	return len(d.Opened) > 0 || len(d.Closed) > 0
}

func index(states []scanner.PortState) map[int]scanner.PortState {
	m := make(map[int]scanner.PortState, len(states))
	for _, s := range states {
		m[s.Port] = s
	}
	return m
}
