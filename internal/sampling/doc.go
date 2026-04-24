// Package sampling implements probabilistic event sampling for the
// portwatch alert pipeline.
//
// # Overview
//
// When a monitored port flaps rapidly, the volume of alerts can become
// overwhelming. The Sampler reduces noise by forwarding only a
// configurable fraction of events while still providing a statistical
// signal that something is wrong.
//
// # Usage
//
//	s := sampling.New(0.1) // forward ~10 % of events
//	if s.Allow(port) {
//		notifier.Notify(ctx, event)
//	}
//
// A rate of 1.0 (the default for invalid inputs) disables sampling and
// forwards every event, preserving the existing behaviour of the pipeline.
package sampling
