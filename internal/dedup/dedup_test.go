package dedup

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllow_FirstCallPermitted(t *testing.T) {
	d := NewWithClock(10*time.Second, fixedClock(epoch))
	if !d.Allow(8080, "opened") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_DuplicateWithinWindowBlocked(t *testing.T) {
	d := NewWithClock(10*time.Second, fixedClock(epoch))
	d.Allow(8080, "opened")
	if d.Allow(8080, "opened") {
		t.Fatal("expected duplicate within window to be blocked")
	}
}

func TestAllow_AfterWindowPermitted(t *testing.T) {
	now := epoch
	d := NewWithClock(10*time.Second, func() time.Time { return now })
	d.Allow(8080, "opened")
	now = epoch.Add(11 * time.Second)
	if !d.Allow(8080, "opened") {
		t.Fatal("expected call after window expiry to be allowed")
	}
}

func TestAllow_DifferentDirectionsAreIndependent(t *testing.T) {
	d := NewWithClock(10*time.Second, fixedClock(epoch))
	d.Allow(8080, "opened")
	if !d.Allow(8080, "closed") {
		t.Fatal("expected different direction to be allowed independently")
	}
}

func TestAllow_DifferentPortsAreIndependent(t *testing.T) {
	d := NewWithClock(10*time.Second, fixedClock(epoch))
	d.Allow(8080, "opened")
	if !d.Allow(9090, "opened") {
		t.Fatal("expected different port to be allowed independently")
	}
}

func TestReset_AllowsImmediateRepeat(t *testing.T) {
	d := NewWithClock(10*time.Second, fixedClock(epoch))
	d.Allow(8080, "opened")
	d.Reset(8080, "opened")
	if !d.Allow(8080, "opened") {
		t.Fatal("expected allow after reset")
	}
}

func TestPurge_RemovesExpiredEntries(t *testing.T) {
	now := epoch
	d := NewWithClock(10*time.Second, func() time.Time { return now })
	d.Allow(8080, "opened")
	d.Allow(9090, "opened")

	now = epoch.Add(11 * time.Second)
	d.Purge()

	if len(d.seen) != 0 {
		t.Fatalf("expected empty map after purge, got %d entries", len(d.seen))
	}
}

func TestPurge_KeepsActiveEntries(t *testing.T) {
	now := epoch
	d := NewWithClock(10*time.Second, func() time.Time { return now })
	d.Allow(8080, "opened")

	now = epoch.Add(5 * time.Second)
	d.Purge()

	if len(d.seen) != 1 {
		t.Fatalf("expected 1 active entry after purge, got %d", len(d.seen))
	}
}
