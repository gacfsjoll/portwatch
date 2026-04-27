package window_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/window"
)

// fixedClock returns a clock whose current time can be advanced manually.
func fixedClock(initial time.Time) (clock func() time.Time, advance func(time.Duration)) {
	now := initial
	return func() time.Time { return now },
		func(d time.Duration) { now = now.Add(d) }
}

func TestTotal_EmptyCounter(t *testing.T) {
	c := window.New(time.Minute)
	if got := c.Total(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestAdd_CountsWithinWindow(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	clk, _ := fixedClock(base)
	c := window.NewWithClock(time.Minute, clk)
	c.Add(3)
	c.Add(5)
	if got := c.Total(); got != 8 {
		t.Fatalf("expected 8, got %d", got)
	}
}

func TestAdd_EvictsExpiredEntries(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	clk, advance := fixedClock(base)
	c := window.NewWithClock(time.Minute, clk)

	c.Add(10) // recorded at t=0
	advance(90 * time.Second) // move past the 60 s window
	c.Add(2)  // recorded at t=90s

	if got := c.Total(); got != 2 {
		t.Fatalf("expected 2 (old entry evicted), got %d", got)
	}
}

func TestAdd_EntryExactlyAtWindowBoundaryEvicted(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	clk, advance := fixedClock(base)
	c := window.NewWithClock(time.Minute, clk)

	c.Add(7)
	advance(time.Minute) // exactly one window later — entry is now stale
	c.Add(1)

	if got := c.Total(); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestReset_ClearsAll(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	clk, _ := fixedClock(base)
	c := window.NewWithClock(time.Minute, clk)
	c.Add(4)
	c.Add(6)
	c.Reset()
	if got := c.Total(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestTotal_MultipleWindowsAccumulate(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	clk, advance := fixedClock(base)
	c := window.NewWithClock(30*time.Second, clk)

	c.Add(1)
	advance(10 * time.Second)
	c.Add(2)
	advance(10 * time.Second)
	c.Add(3)
	// all three within 30 s window
	if got := c.Total(); got != 6 {
		t.Fatalf("expected 6, got %d", got)
	}

	advance(15 * time.Second) // first entry (t=0) is now 25 s old — still inside
	if got := c.Total(); got != 6 {
		t.Fatalf("expected 6, got %d", got)
	}

	advance(10 * time.Second) // first entry is now 35 s old — evicted
	if got := c.Total(); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}
