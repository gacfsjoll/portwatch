package throttle_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) throttle.Clock {
	return func() time.Time { return t }
}

func TestAllow_FirstCallPermitted(t *testing.T) {
	th := throttle.NewWithClock(time.Minute, 3, fixedClock(epoch))
	if !th.Allow(8080) {
		t.Fatal("expected first alert to be allowed")
	}
}

func TestAllow_BlockedAfterMaxHits(t *testing.T) {
	th := throttle.NewWithClock(time.Minute, 2, fixedClock(epoch))
	if !th.Allow(9000) {
		t.Fatal("hit 1 should be allowed")
	}
	if !th.Allow(9000) {
		t.Fatal("hit 2 should be allowed")
	}
	if th.Allow(9000) {
		t.Fatal("hit 3 should be blocked")
	}
}

func TestAllow_ResetsAfterWindow(t *testing.T) {
	now := epoch
	clock := func() time.Time { return now }
	th := throttle.NewWithClock(time.Minute, 1, clock)

	th.Allow(443)
	if th.Allow(443) {
		t.Fatal("second call within window should be blocked")
	}

	now = epoch.Add(2 * time.Minute)
	if !th.Allow(443) {
		t.Fatal("call after window should be allowed")
	}
}

func TestAllow_DifferentPortsIndependent(t *testing.T) {
	th := throttle.NewWithClock(time.Minute, 1, fixedClock(epoch))
	th.Allow(80)
	if !th.Allow(443) {
		t.Fatal("different port should not be affected")
	}
}

func TestReset_ClearsHits(t *testing.T) {
	th := throttle.NewWithClock(time.Minute, 1, fixedClock(epoch))
	th.Allow(8080)
	th.Reset(8080)
	if !th.Allow(8080) {
		t.Fatal("after reset, port should be allowed again")
	}
}

func TestStats_ReturnsCount(t *testing.T) {
	th := throttle.NewWithClock(time.Minute, 5, fixedClock(epoch))
	th.Allow(3000)
	th.Allow(3000)
	if got := th.Stats(3000); got != 2 {
		t.Fatalf("expected 2 hits, got %d", got)
	}
}
