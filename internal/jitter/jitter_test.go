package jitter_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/user/portwatch/internal/jitter"
)

// deterministicSource always returns the same value so tests are stable.
type deterministicSource struct{ val int64 }

func (d *deterministicSource) Int63() int64 { return d.val }
func (d *deterministicSource) Seed(_ int64) {}

func withSource(val int64, fn func()) {
	orig := jitter.Source
	jitter.Source = rand.NewSource(val)
	defer func() { jitter.Source = orig }()
	fn()
}

func TestApply_ZeroJitter_ReturnsBase(t *testing.T) {
	withSource(42, func() {
		got := jitter.Apply(10*time.Second, 0)
		if got != 10*time.Second {
			t.Fatalf("expected 10s, got %v", got)
		}
	})
}

func TestApply_NegativeJitter_ReturnsBase(t *testing.T) {
	withSource(42, func() {
		got := jitter.Apply(10*time.Second, -1*time.Second)
		if got != 10*time.Second {
			t.Fatalf("expected 10s, got %v", got)
		}
	})
}

func TestApply_ResultWithinBounds(t *testing.T) {
	base := 10 * time.Second
	maxJ := 2 * time.Second

	for i := 0; i < 200; i++ {
		got := jitter.Apply(base, maxJ)
		if got < base-maxJ || got > base+maxJ {
			t.Fatalf("iteration %d: %v out of [%v, %v]", i, got, base-maxJ, base+maxJ)
		}
	}
}

func TestApply_JitterLargerThanBase_Clamped(t *testing.T) {
	base := 5 * time.Second
	maxJ := 20 * time.Second // will be clamped to base

	for i := 0; i < 100; i++ {
		got := jitter.Apply(base, maxJ)
		if got < 0 {
			t.Fatalf("iteration %d: got negative duration %v", i, got)
		}
	}
}

func TestPercent_ZeroPct_ReturnsBase(t *testing.T) {
	got := jitter.Percent(10*time.Second, 0)
	if got != 10*time.Second {
		t.Fatalf("expected 10s, got %v", got)
	}
}

func TestPercent_ResultWithinBounds(t *testing.T) {
	base := 30 * time.Second
	pct := 10 // ±3 s

	for i := 0; i < 200; i++ {
		got := jitter.Percent(base, pct)
		lo := 27 * time.Second
		hi := 33 * time.Second
		if got < lo || got > hi {
			t.Fatalf("iteration %d: %v out of [%v, %v]", i, got, lo, hi)
		}
	}
}

func TestPercent_Over100_ClampsTo100(t *testing.T) {
	base := 4 * time.Second
	for i := 0; i < 100; i++ {
		got := jitter.Percent(base, 150)
		if got < 0 {
			t.Fatalf("iteration %d: negative duration %v", i, got)
		}
	}
}
