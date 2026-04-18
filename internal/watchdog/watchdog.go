// Package watchdog provides a self-health check mechanism that verifies
// the portwatch daemon is running and scanning as expected.
package watchdog

import (
	"fmt"
	"os"
	"time"
)

// Status represents the current health of the daemon.
type Status struct {
	Healthy   bool
	LastScan  time.Time
	ScanCount int64
	Message   string
}

// Watchdog monitors the daemon's internal health.
type Watchdog struct {
	maxStaleness time.Duration
	lastScan     time.Time
	scanCount    int64
	pidFile      string
}

// New creates a Watchdog with the given staleness threshold and PID file path.
func New(maxStaleness time.Duration, pidFile string) *Watchdog {
	return &Watchdog{
		maxStaleness: maxStaleness,
		pidFile:      pidFile,
	}
}

// RecordScan updates the last scan timestamp and increments the counter.
func (w *Watchdog) RecordScan() {
	w.lastScan = time.Now()
	w.scanCount++
}

// Check returns the current health status of the daemon.
func (w *Watchdog) Check() Status {
	if w.lastScan.IsZero() {
		return Status{Healthy: false, Message: "no scan has been recorded yet"}
	}
	staleness := time.Since(w.lastScan)
	if staleness > w.maxStaleness {
		return Status{
			Healthy:   false,
			LastScan:  w.lastScan,
			ScanCount: w.scanCount,
			Message:   fmt.Sprintf("last scan was %s ago (threshold: %s)", staleness.Round(time.Second), w.maxStaleness),
		}
	}
	return Status{
		Healthy:   true,
		LastScan:  w.lastScan,
		ScanCount: w.scanCount,
		Message:   "ok",
	}
}

// WritePID writes the current process PID to the configured pid file.
func (w *Watchdog) WritePID() error {
	if w.pidFile == "" {
		return nil
	}
	data := fmt.Sprintf("%d\n", os.Getpid())
	return os.WriteFile(w.pidFile, []byte(data), 0o644)
}

// RemovePID removes the PID file on shutdown.
func (w *Watchdog) RemovePID() {
	if w.pidFile != "" {
		_ = os.Remove(w.pidFile)
	}
}
