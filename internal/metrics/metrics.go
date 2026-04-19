// Package metrics tracks runtime counters for portwatch scans and alerts.
package metrics

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Counters holds cumulative metrics collected during a portwatch session.
type Counters struct {
	mu           sync.Mutex
	ScansTotal   int
	AlertsTotal  int
	OpenedTotal  int
	ClosedTotal  int
	LastScanTime time.Time
}

// Recorder records scan and alert metrics.
type Recorder struct {
	counters Counters
	clock    func() time.Time
}

// New returns a Recorder using the real clock.
func New() *Recorder {
	return &Recorder{clock: time.Now}
}

// NewWithClock returns a Recorder with an injectable clock (for testing).
func NewWithClock(clock func() time.Time) *Recorder {
	return &Recorder{clock: clock}
}

// RecordScan increments the scan counter and updates the last scan timestamp.
func (r *Recorder) RecordScan() {
	r.counters.mu.Lock()
	defer r.counters.mu.Unlock()
	r.counters.ScansTotal++
	r.counters.LastScanTime = r.clock()
}

// RecordAlert increments the alert and directional counters.
func (r *Recorder) RecordAlert(opened, closed int) {
	r.counters.mu.Lock()
	defer r.counters.mu.Unlock()
	r.counters.AlertsTotal++
	r.counters.OpenedTotal += opened
	r.counters.ClosedTotal += closed
}

// Snapshot returns a copy of the current counters.
func (r *Recorder) Snapshot() Counters {
	r.counters.mu.Lock()
	defer r.counters.mu.Unlock()
	return Counters{
		ScansTotal:   r.counters.ScansTotal,
		AlertsTotal:  r.counters.AlertsTotal,
		OpenedTotal:  r.counters.OpenedTotal,
		ClosedTotal:  r.counters.ClosedTotal,
		LastScanTime: r.counters.LastScanTime,
	}
}

// Print writes a human-readable summary to w (defaults to os.Stdout).
func (r *Recorder) Print(w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	s := r.Snapshot()
	fmt.Fprintf(w, "scans=%d alerts=%d opened=%d closed=%d last_scan=%s\n",
		s.ScansTotal, s.AlertsTotal, s.OpenedTotal, s.ClosedTotal,
		s.LastScanTime.Format(time.RFC3339))
}
