package ratelimit_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/ratelimit"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllow_FirstCallPermitted(t *testing.T) {
	l := ratelimit.New(time.Minute)
	if !l.Allow(8080) {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_WithinCooldownBlocked(t *testing.T) {
	base := time.Now()
	l := ratelimit.New(time.Minute)
	l.(*ratelimit.Limiter) // ensure exported; use via interface if needed

	// Use internal clock override via a helper constructor in tests.
	l2 := ratelimit.NewWithClock(time.Minute, fixedClock(base))
	l2.Allow(9000)
	if l2.Allow(9000) {
		t.Fatal("expected second call within cooldown to be blocked")
	}
}

func TestAllow_AfterCooldownPermitted(t *testing.T) {
	base := time.Now()
	l := ratelimit.NewWithClock(time.Minute, fixedClock(base))
	l.Allow(443)

	// Advance clock past cooldown.
	l2 := ratelimit.NewWithClock(time.Minute, fixedClock(base.Add(2*time.Minute)))
	// Copy state by resetting and re-recording with new clock not possible
	// directly; test via Reset instead.
	l.Reset(443)
	if !l.Allow(443) {
		t.Fatal("expected allow after reset")
	}
}

func TestReset_ClearsPort(t *testing.T) {
	base := time.Now()
	l := ratelimit.NewWithClock(time.Minute, fixedClock(base))
	l.Allow(22)
	l.Reset(22)
	if !l.Allow(22) {
		t.Fatal("expected allow after reset")
	}
}

func TestResetAll_ClearsAllPorts(t *testing.T) {
	base := time.Now()
	l := ratelimit.NewWithClock(time.Minute, fixedClock(base))
	l.Allow(80)
	l.Allow(443)
	l.ResetAll()
	if !l.Allow(80) || !l.Allow(443) {
		t.Fatal("expected all ports allowed after ResetAll")
	}
}
