package debounce_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/debounce"
)

var (
	t0 = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
)

func fixedClock(t time.Time) debounce.Clock {
	return func() time.Time { return t }
}

func TestAllow_FirstCallPermitted(t *testing.T) {
	d := debounce.NewWithClock(5*time.Second, fixedClock(t0))
	if !d.Allow("port:8080:opened") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_WithinWindowBlocked(t *testing.T) {
	d := debounce.NewWithClock(5*time.Second, fixedClock(t0))
	d.Allow("port:8080:opened")

	// same instant — within window
	if d.Allow("port:8080:opened") {
		t.Fatal("expected second call within window to be blocked")
	}
}

func TestAllow_AfterWindowPermitted(t *testing.T) {
	current := t0
	clock := func() time.Time { return current }
	d := debounce.NewWithClock(5*time.Second, clock)

	d.Allow("port:8080:opened")
	current = t0.Add(6 * time.Second)

	if !d.Allow("port:8080:opened") {
		t.Fatal("expected call after window to be allowed")
	}
}

func TestAllow_DifferentKeysIndependent(t *testing.T) {
	d := debounce.NewWithClock(5*time.Second, fixedClock(t0))
	d.Allow("port:8080:opened")

	if !d.Allow("port:9090:opened") {
		t.Fatal("expected different key to be allowed")
	}
}

func TestReset_ClearsKey(t *testing.T) {
	d := debounce.NewWithClock(5*time.Second, fixedClock(t0))
	d.Allow("port:8080:opened")
	d.Reset("port:8080:opened")

	if !d.Allow("port:8080:opened") {
		t.Fatal("expected allow after reset")
	}
}

func TestFlush_ClearsAll(t *testing.T) {
	d := debounce.NewWithClock(5*time.Second, fixedClock(t0))
	d.Allow("port:8080:opened")
	d.Allow("port:9090:opened")
	d.Flush()

	if !d.Allow("port:8080:opened") {
		t.Fatal("expected allow after flush")
	}
	if !d.Allow("port:9090:opened") {
		t.Fatal("expected allow after flush")
	}
}
