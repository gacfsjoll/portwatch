// Package cooldown implements a per-port quiet-period tracker designed to
// suppress alert storms caused by flapping ports.
//
// A flapping port is one that transitions between open and closed states
// repeatedly in a short window. Without cooldown, each transition would
// generate an alert, flooding the operator with noise.
//
// Usage:
//
//	tr := cooldown.New(30 * time.Second)
//
//	if tr.Allow(port) {
//		// send alert
//	}
//
// The quiet period begins when an alert is first allowed through. Subsequent
// calls for the same port within the period return false. Once the period
// elapses the next call is permitted and the timer resets.
//
// Reset can be used to clear the quiet period for a specific port, for
// example when an operator explicitly acknowledges the port.
package cooldown
