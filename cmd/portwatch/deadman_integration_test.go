package main

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/deadman"
)

// TestDeadman_IntegrationWithMonitorLoop simulates a monitor loop that calls
// Reset after each scan and verifies the switch never trips during normal
// operation, then confirms it trips once the loop stalls.
func TestDeadman_IntegrationWithMonitorLoop(t *testing.T) {
	const deadline = 150 * time.Millisecond
	const scanInterval = 40 * time.Millisecond

	var alerts int32
	sw := deadman.New(deadline, func(_ time.Time, _ time.Duration) {
		atomic.AddInt32(&alerts, 1)
	})

	// Phase 1: healthy loop — reset regularly, no alerts expected.
	stop := make(chan struct{})
	go func() {
		for {
			select {
			case <-stop:
				return
			case <-time.After(scanInterval):
				sw.Reset()
			}
		}
	}()

	time.Sleep(400 * time.Millisecond)
	close(stop) // stop the healthy loop

	if n := atomic.LoadInt32(&alerts); n != 0 {
		t.Fatalf("phase 1: expected 0 alerts during healthy loop, got %d", n)
	}

	// Phase 2: loop has stopped — switch should trip.
	time.Sleep(400 * time.Millisecond)
	sw.Stop()

	if n := atomic.LoadInt32(&alerts); n == 0 {
		t.Fatal("phase 2: expected at least one alert after loop stalled, got 0")
	}
}
