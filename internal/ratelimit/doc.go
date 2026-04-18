// Package ratelimit implements per-port alert rate limiting for portwatch.
//
// When the monitor detects rapid repeated open/close events on the same port
// (e.g. due to short-lived connections or a flapping service), ratelimit
// prevents the alert notifier from being flooded with duplicate notifications.
//
// Usage:
//
//	limiter := ratelimit.New(5 * time.Minute)
//
//	if limiter.Allow(port) {
//		notifier.Notify(event)
//	}
//
// The cooldown window is configurable and defaults to the value set in the
// portwatch configuration file under the alert.cooldown key.
package ratelimit
