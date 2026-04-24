// Package jitter provides randomised interval spreading to prevent
// thundering-herd problems when multiple portwatch instances start
// simultaneously or when scan loops are tightly coupled to wall-clock
// boundaries.
package jitter

import (
	"math/rand"
	"sync"
	"time"
)

// Source is the random-number source used by Apply. It is a package-level
// variable so tests can swap it out for a deterministic implementation.
var Source rand.Source = rand.NewSource(time.Now().UnixNano())

var mu sync.Mutex

// Apply returns base ± a random fraction of maxJitter.
// maxJitter is clamped to base so the result is always positive.
//
//	Apply(10s, 2s)  →  somewhere in [8s, 12s]
func Apply(base, maxJitter time.Duration) time.Duration {
	if maxJitter <= 0 {
		return base
	}
	if maxJitter > base {
		maxJitter = base
	}

	mu.Lock()
	r := rand.New(Source) //nolint:gosec // non-crypto use
	offset := time.Duration(r.Int63n(int64(maxJitter)*2+1)) - maxJitter
	mu.Unlock()

	return base + offset
}

// Percent returns Apply(base, base*pct/100), where pct is a percentage
// between 0 and 100. Values outside that range are clamped.
//
//	Percent(10s, 20)  →  somewhere in [8s, 12s]
func Percent(base time.Duration, pct int) time.Duration {
	if pct <= 0 {
		return base
	}
	if pct > 100 {
		pct = 100
	}
	maxJitter := time.Duration(float64(base) * float64(pct) / 100.0)
	return Apply(base, maxJitter)
}
