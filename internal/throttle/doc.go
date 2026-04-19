// Package throttle provides a rolling-window rate limiter for port change
// alerts. It is distinct from ratelimit (which enforces a per-port cooldown
// between consecutive alerts) in that throttle counts total alert firings
// within a configurable time window and suppresses any that exceed a
// configurable maximum. This prevents alert storms when many ports open or
// close in rapid succession.
//
// Basic usage:
//
//	th := throttle.New(time.Minute, 5)
//	if th.Allow(port) {
//		notifier.Notify(event)
//	}
package throttle
