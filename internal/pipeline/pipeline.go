// Package pipeline wires together the alert processing stages — filtering,
// suppression, acknowledgement, debouncing, rate-limiting, throttling, rollup,
// and circuit-breaking — into a single ordered chain that every port-change
// event passes through before a notifier is called.
package pipeline

import (
	"context"
	"log"

	"github.com/user/portwatch/internal/acknowledge"
	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/circuitbreaker"
	"github.com/user/portwatch/internal/debounce"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/ratelimit"
	"github.com/user/portwatch/internal/rollup"
	"github.com/user/portwatch/internal/suppress"
	"github.com/user/portwatch/internal/throttle"
)

// Stage is a single processing step.  It returns false to drop the event.
type Stage func(evt alert.Event) bool

// Pipeline runs an event through an ordered sequence of stages and, if every
// stage accepts it, forwards it to the wrapped Notifier.
type Pipeline struct {
	stages   []Stage
	sink     alert.Notifier
}

// New constructs a Pipeline from the provided stages and terminal notifier.
func New(sink alert.Notifier, stages ...Stage) *Pipeline {
	return &Pipeline{stages: stages, sink: sink}
}

// Notify implements alert.Notifier.  Each stage is evaluated in order; the
// first stage that returns false causes the event to be silently dropped.
func (p *Pipeline) Notify(ctx context.Context, evt alert.Event) error {
	for _, s := range p.stages {
		if !s(evt) {
			return nil
		}
	}
	return p.sink.Notify(ctx, evt)
}

// ---------------------------------------------------------------------------
// Stage constructors
// ---------------------------------------------------------------------------

// FilterStage drops events whose port is not allowed by the filter.
func FilterStage(f *filter.Filter) Stage {
	return func(evt alert.Event) bool {
		return f.Allow(evt.Port)
	}
}

// SuppressStage drops events whose port is currently suppressed.
func SuppressStage(s *suppress.Suppressor) Stage {
	return func(evt alert.Event) bool {
		return !s.IsSuppressed(evt.Port)
	}
}

// AcknowledgeStage drops events for ports that have been acknowledged.
func AcknowledgeStage(store *acknowledge.Store) Stage {
	return func(evt alert.Event) bool {
		return !store.IsAcknowledged(evt.Port)
	}
}

// DebounceStage drops repeated events for the same port within the debounce
// window to avoid alert storms during rapid state flapping.
func DebounceStage(d *debounce.Debouncer) Stage {
	return func(evt alert.Event) bool {
		return d.Allow(evt.Port)
	}
}

// RateLimitStage enforces a per-port cooldown between successive alerts.
func RateLimitStage(r *ratelimit.Limiter) Stage {
	return func(evt alert.Event) bool {
		return r.Allow(evt.Port)
	}
}

// ThrottleStage caps the number of alerts per port within a sliding window.
func ThrottleStage(t *throttle.Throttle) Stage {
	return func(evt alert.Event) bool {
		return t.Allow(evt.Port)
	}
}

// RollupStage batches events through the rollup buffer.  Because rollup has
// its own internal flush goroutine the stage always returns false — the rollup
// buffer is responsible for forwarding batched events to the downstream sink
// separately.  Callers must start the rollup flush loop before using this.
func RollupStage(r *rollup.Buffer) Stage {
	return func(evt alert.Event) bool {
		r.Add(evt)
		return false // rollup owns delivery
	}
}

// CircuitBreakerStage opens the circuit after repeated notifier failures,
// preventing alert floods when a downstream system is unavailable.
func CircuitBreakerStage(cb *circuitbreaker.CircuitBreaker) Stage {
	return func(evt alert.Event) bool {
		if !cb.Allow() {
			log.Printf("[pipeline] circuit open — dropping alert for port %d", evt.Port)
			return false
		}
		return true
	}
}
