// Package circuitbreaker provides a thread-safe circuit breaker for protecting
// downstream notifiers (webhooks, log sinks, etc.) from cascading failures.
//
// # Overview
//
// A Breaker tracks consecutive failures against a configured threshold. Once
// the threshold is reached the circuit transitions to the open state and all
// subsequent Allow calls return ErrCircuitOpen until the recovery window has
// elapsed.
//
// After the recovery window the circuit moves back to closed on the next Allow
// call, allowing a single probe request through. If that probe succeeds (the
// caller invokes RecordSuccess) the breaker remains closed. If it fails the
// circuit opens again immediately.
//
// # Usage
//
//	b := circuitbreaker.New(5, 30*time.Second)
//
//	if err := b.Allow(); err != nil {
//		// skip notification — circuit is open
//		return err
//	}
//	if err := notifier.Notify(ctx, event); err != nil {
//		b.RecordFailure()
//		return err
//	}
//	b.RecordSuccess()
package circuitbreaker
