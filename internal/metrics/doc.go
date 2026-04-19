// Package metrics provides lightweight runtime counters for the portwatch
// daemon. It tracks the total number of port scans performed, alert events
// fired, and the directional breakdown of opened vs closed ports observed
// during the current process lifetime.
//
// Usage:
//
//	rec := metrics.New()
//	rec.RecordScan()
//	rec.RecordAlert(opened, closed)
//	s := rec.Snapshot()   // thread-safe copy
//	rec.Print(os.Stdout)  // human-readable summary
//
// All methods are safe for concurrent use.
package metrics
