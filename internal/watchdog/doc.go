// Package watchdog provides daemon self-health monitoring for portwatch.
//
// It tracks when the last port scan occurred and exposes a Check method
// that returns whether the daemon is operating within expected parameters.
//
// # Usage
//
//	w := watchdog.New(2*time.Minute, "/var/run/portwatch.pid")
//	w.WritePID()
//	defer w.RemovePID()
//
//	// after each scan cycle:
//	w.RecordScan()
//
//	// to inspect health:
//	status := w.Check()
//	if !status.Healthy {
//		log.Println("watchdog:", status.Message)
//	}
package watchdog
