// Package acknowledge provides persistent tracking of operator-acknowledged
// open ports within portwatch.
//
// When an unexpected port is detected, the operator may choose to acknowledge
// it — indicating the port is known and expected. Acknowledged ports are
// persisted to disk and will not trigger further alerts until their
// acknowledgement is revoked or the port closes and reopens.
//
// Usage:
//
//	store := acknowledge.NewStore("/var/lib/portwatch/acks.json")
//	if err := store.Load(); err != nil { ... }
//
//	// Silence alerts for port 8080
//	store.Acknowledge(8080)
//
//	// Check before alerting
//	if !store.IsAcknowledged(port) {
//	    notifier.Notify(event)
//	}
package acknowledge
