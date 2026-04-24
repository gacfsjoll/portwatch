// Package rollup provides a time-windowed event batcher for port-change
// alerts.
//
// When portwatch detects several port changes in rapid succession — for
// example when a service restarts and briefly closes then re-opens a set of
// ports — emitting one alert per change can be overwhelming.  The Roller
// type collects events that arrive within a configurable quiet window and
// delivers them together as a single batch to a Flusher callback.
//
// Typical usage:
//
//	roller := rollup.New(500*time.Millisecond, func(events []alert.Event) {
//		for _, e := range events {
//			notifier.Notify(ctx, e)
//		}
//	})
//
//	// Later, on each detected change:
//	roller.Add(event)
//
//	// On daemon shutdown:
//	roller.Flush()
package rollup
