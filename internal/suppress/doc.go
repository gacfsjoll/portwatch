// Package suppress implements a time-based port alert suppression list.
//
// During planned maintenance or known service restarts, individual ports
// can be suppressed for a configurable duration so that the monitor does
// not emit spurious opened/closed alerts.
//
// Usage:
//
//	list := suppress.New()
//	list.Suppress(8080, 30*time.Minute, "deploy in progress")
//
//	// later, inside alert pipeline:
//	if list.IsSuppressed(event.Port) {
//	    return nil // skip notification
//	}
//
// Expired entries can be cleaned up periodically via Expire.
package suppress
