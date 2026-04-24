package cooldown_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/cooldown"
)

// fixedClock returns a clock whose value can be advanced manually.
func fixedClock(initial time.Time) (func() time.Time, func(time.Duration)) {
	current := initial
	get := func() time.Time { return current }
	advance := func(d time.Duration) { current = current.Add(d) }
	return get, advance
}

func TestAllow_FirstCallPermitted(t *testing.T) {
	clk, _ := fixedClock(time.Unix(1000, 0))
	tr := cooldown.NewWithClock(5*time.Second, clk)

	if !tr.Allow(8080) {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_WithinPeriodBlocked(t *testing.T) {
	clk, advance := fixedClock(time.Unix(1000, 0))
	tr := cooldown.NewWithClock(5*time.Second, clk)

	tr.Allow(8080)
	advance(2 * time.Second)

	if tr.Allow(8080) {
		t.Fatal("expected call within quiet period to be blocked")
	}
}

func TestAllow_AfterPeriodPermitted(t *testing.T) {
	clk, advance := fixedClock(time.Unix(1000, 0))
	tr := cooldown.NewWithClock(5*time.Second, clk)

	tr.Allow(8080)
	advance(6 * time.Second)

	if !tr.Allow(8080) {
		t.Fatal("expected call after quiet period to be allowed")
	}
}

func TestAllow_DifferentPortsIndependent(t *testing.T) {
	clk, _ := fixedClock(time.Unix(1000, 0))
	tr := cooldown.NewWithClock(5*time.Second, clk)

	tr.Allow(8080)

	if !tr.Allow(9090) {
		t.Fatal("expected different port to be allowed independently")
	}
}

func TestReset_AllowsImmediately(t *testing.T) {
	clk, _ := fixedClock(time.Unix(1000, 0))
	tr := cooldown.NewWithClock(5*time.Second, clk)

	tr.Allow(8080)
	tr.Reset(8080)

	if !tr.Allow(8080) {
		t.Fatal("expected reset port to be allowed immediately")
	}
}

func TestActive_CountsPortsInPeriod(t *testing.T) {
	clk, advance := fixedClock(time.Unix(1000, 0))
	tr := cooldown.NewWithClock(10*time.Second, clk)

	tr.Allow(80)
	tr.Allow(443)
	tr.Allow(8080)

	if got := tr.Active(); got != 3 {
		t.Fatalf("expected 3 active, got %d", got)
	}

	advance(11 * time.Second)

	if got := tr.Active(); got != 0 {
		t.Fatalf("expected 0 active after period, got %d", got)
	}
}
