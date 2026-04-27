package limiter

import (
	"testing"
	"time"
)

// fixedClock returns a Clock that always returns t.
func fixedClock(t time.Time) Clock {
	return func() time.Time { return t }
}

func TestAllow_FirstCallPermitted(t *testing.T) {
	l := NewWithClock(3, time.Minute, fixedClock(time.Now()))
	if !l.Allow(8080) {
		t.Fatal("expected first call to be permitted")
	}
}

func TestAllow_BlockedAfterMaxHits(t *testing.T) {
	now := time.Now()
	l := NewWithClock(2, time.Minute, fixedClock(now))

	if !l.Allow(9000) {
		t.Fatal("call 1 should be permitted")
	}
	if !l.Allow(9000) {
		t.Fatal("call 2 should be permitted")
	}
	if l.Allow(9000) {
		t.Fatal("call 3 should be blocked (over max)")
	}
}

func TestAllow_ResetsAfterWindow(t *testing.T) {
	base := time.Now()
	clock := base
	l := NewWithClock(1, time.Minute, func() time.Time { return clock })

	if !l.Allow(443) {
		t.Fatal("first call should be permitted")
	}
	if l.Allow(443) {
		t.Fatal("second call within window should be blocked")
	}

	// Advance past the window.
	clock = base.Add(2 * time.Minute)
	if !l.Allow(443) {
		t.Fatal("call after window expiry should be permitted")
	}
}

func TestAllow_DifferentPortsIndependent(t *testing.T) {
	l := NewWithClock(1, time.Minute, fixedClock(time.Now()))

	l.Allow(80)
	if l.Allow(80) {
		t.Fatal("port 80 second call should be blocked")
	}
	if !l.Allow(443) {
		t.Fatal("port 443 first call should be permitted independently")
	}
}

func TestReset_ClearsBucket(t *testing.T) {
	now := time.Now()
	l := NewWithClock(1, time.Minute, fixedClock(now))

	l.Allow(22)
	if l.Allow(22) {
		t.Fatal("should be blocked before reset")
	}
	l.Reset(22)
	if !l.Allow(22) {
		t.Fatal("should be permitted after reset")
	}
}

func TestStats_ReturnsCountAndSince(t *testing.T) {
	now := time.Now()
	l := NewWithClock(5, time.Minute, fixedClock(now))

	count, since := l.Stats(3306)
	if count != 0 || !since.IsZero() {
		t.Fatal("expected zero stats for unseen port")
	}

	l.Allow(3306)
	l.Allow(3306)
	count, since = l.Stats(3306)
	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
	if !since.Equal(now) {
		t.Fatalf("expected since=%v, got %v", now, since)
	}
}
