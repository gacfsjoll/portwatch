// Package deadman provides a dead-man's switch for the portwatch monitor.
//
// A dead-man's switch is a safety mechanism that fires a notification when
// it has NOT been explicitly reset within a configured deadline. In
// portwatch this is used to detect silent scan failures: if the monitor
// goroutine stalls or panics without the daemon exiting, the switch will
// alert the operator that scans have stopped.
//
// # Usage
//
//	sw := deadman.New(2*time.Minute, func(last time.Time, elapsed time.Duration) {
//	    log.Printf("[WARN] no scan completed in %s (last seen %s)", elapsed, last)
//	})
//	defer sw.Stop()
//
//	// Inside the scan loop:
//	sw.Reset()  // called after every successful scan
package deadman
