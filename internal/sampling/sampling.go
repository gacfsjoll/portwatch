// Package sampling provides probabilistic event sampling for portwatch.
// It allows reducing alert volume by forwarding only a statistical
// fraction of events, useful when a port is flapping frequently.
package sampling

import (
	"math/rand"
	"sync"
	"time"
)

// Sampler decides whether an event for a given port should be forwarded
// based on a configured sample rate in the range (0.0, 1.0].
type Sampler struct {
	mu   sync.Mutex
	rate float64
	rng  *rand.Rand
}

// New returns a Sampler that forwards events with the given probability.
// A rate of 1.0 forwards every event; 0.5 forwards roughly half.
// Values outside (0, 1] are clamped to 1.0.
func New(rate float64) *Sampler {
	if rate <= 0 || rate > 1 {
		rate = 1.0
	}
	return &Sampler{
		rate: rate,
		rng:  rand.New(rand.NewSource(time.Now().UnixNano())), //nolint:gosec
	}
}

// Allow reports whether the event for port should be forwarded.
// It is safe for concurrent use.
func (s *Sampler) Allow(port int) bool {
	if s.rate >= 1.0 {
		return true
	}
	s.mu.Lock()
	v := s.rng.Float64()
	s.mu.Unlock()
	_ = port // port reserved for per-port rate maps in future
	return v < s.rate
}

// Rate returns the configured sample rate.
func (s *Sampler) Rate() float64 {
	return s.rate
}
