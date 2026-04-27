// Package replay re-emits historical alert events through a notifier pipeline.
//
// It is useful in two scenarios:
//
//  1. Testing alert routing — you can replay a recorded history file through a
//     new notifier configuration to verify that events are delivered correctly.
//
//  2. Recovery — if the daemon was restarted and some notifications were missed,
//     replay can re-dispatch events from the last N hours through the configured
//     notifier without re-scanning ports.
//
// # Usage
//
//	opts := replay.Options{Since: 2 * time.Hour, DryRun: false}
//	r := replay.New(notifier, opts)
//	count, err := r.Run(ctx, recorder)
//
// The Since field limits replay to events newer than the given duration.
// Setting DryRun to true prints events to stdout without dispatching them.
package replay
