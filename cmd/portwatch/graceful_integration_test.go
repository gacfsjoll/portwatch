package main

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/graceful"
)

// TestGraceful_ShutdownPropagates verifies that cancelling the parent
// context (simulating a signal) causes all worker goroutines that
// respect the context to exit within the shutdown window.
func TestGraceful_ShutdownPropagates(t *testing.T) {
	const workers = 4
	var stopped atomic.Int32

	h := graceful.New(2 * time.Second)
	parent, parentCancel := context.WithCancel(context.Background())

	ctx, cancel := h.Wait(parent)
	defer cancel()

	// Spawn workers that exit when ctx is done.
	for i := 0; i < workers; i++ {
		go func() {
			<-ctx.Done()
			stopped.Add(1)
		}()
	}

	// Trigger shutdown.
	parentCancel()

	// All workers should stop within the deadline.
	deadline := time.After(time.Second)
	for {
		if stopped.Load() == int32(workers) {
			return
		}
		select {
		case <-deadline:
			t.Fatalf("only %d/%d workers stopped before deadline", stopped.Load(), workers)
		case <-time.After(10 * time.Millisecond):
		}
	}
}

// TestGraceful_TimeoutBounds confirms that WaitWithTimeout produces a
// context whose deadline is no further than the configured timeout.
func TestGraceful_TimeoutBounds(t *testing.T) {
	timeout := 500 * time.Millisecond
	h := graceful.New(timeout)
	parent, parentCancel := context.WithCancel(context.Background())
	defer parentCancel()

	before := time.Now()
	ctx, cancel := h.WaitWithTimeout(parent)
	defer cancel()

	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected deadline to be set")
	}
	if deadline.Before(before) || deadline.After(before.Add(timeout+50*time.Millisecond)) {
		t.Fatalf("deadline %v out of expected range [%v, %v]", deadline, before, before.Add(timeout))
	}
}
