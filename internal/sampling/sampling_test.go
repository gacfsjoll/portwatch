package sampling_test

import (
	"testing"

	"github.com/user/portwatch/internal/sampling"
)

func TestNew_DefaultsToFullRateForInvalidValues(t *testing.T) {
	for _, rate := range []float64{-1.0, 0.0, 1.5} {
		s := sampling.New(rate)
		if s.Rate() != 1.0 {
			t.Errorf("rate %v: expected 1.0, got %v", rate, s.Rate())
		}
	}
}

func TestNew_StoresValidRate(t *testing.T) {
	s := sampling.New(0.25)
	if s.Rate() != 0.25 {
		t.Fatalf("expected 0.25, got %v", s.Rate())
	}
}

func TestAllow_FullRate_AlwaysPermits(t *testing.T) {
	s := sampling.New(1.0)
	for i := 0; i < 100; i++ {
		if !s.Allow(8080) {
			t.Fatal("expected Allow to return true at rate 1.0")
		}
	}
}

func TestAllow_ZeroRateClamped_AlwaysPermits(t *testing.T) {
	// rate=0 is invalid and clamped to 1.0
	s := sampling.New(0.0)
	if !s.Allow(443) {
		t.Fatal("clamped rate should always permit")
	}
}

func TestAllow_LowRate_RejectsStatistically(t *testing.T) {
	// At rate 0.01 over 1000 trials we expect far fewer than 1000 passes.
	s := sampling.New(0.01)
	allowed := 0
	const trials = 1000
	for i := 0; i < trials; i++ {
		if s.Allow(22) {
			allowed++
		}
	}
	if allowed >= trials {
		t.Errorf("expected sampling to reject some events at rate 0.01, allowed=%d", allowed)
	}
}

func TestAllow_ConcurrentSafe(t *testing.T) {
	s := sampling.New(0.5)
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 50; j++ {
				s.Allow(80)
			}
			done <- struct{}{}
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
}
