// Package baseline manages the trusted snapshot of open ports for portwatch.
//
// # Overview
//
// A baseline represents the set of ports that are considered "expected" on the
// host at a given point in time. Once captured, the monitor compares every
// subsequent scan against the baseline and emits alerts only for ports that
// deviate from it — newly opened ports that are not in the baseline, or
// baseline ports that have unexpectedly closed.
//
// # Usage
//
//	manager := baseline.NewManager("/var/lib/portwatch/baseline.json")
//
//	// Capture the current state:
//	if err := manager.Save(openPorts); err != nil { ... }
//
//	// Later, load and compare:
//	b, err := manager.Load()
//	if errors.Is(err, baseline.ErrNoBaseline) {
//		// prompt user to run 'portwatch capture'
//	}
//	if !b.Contains(port) {
//		// port is unexpected — raise alert
//	}
package baseline
