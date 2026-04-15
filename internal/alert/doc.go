// Package alert defines the Notifier interface and built-in implementations
// used by portwatch to surface port-change events to operators.
//
// # Notifiers
//
// LogNotifier writes human-readable lines to any io.Writer (default: stderr).
//
// MultiNotifier fans a single Event out to an ordered list of Notifiers,
// stopping and returning the first error it encounters.
//
// # Usage
//
//	n := alert.NewLogNotifier(os.Stdout)
//	n.Notify(alert.Event{
//		Timestamp: time.Now(),
//		Level:     alert.LevelAlert,
//		Port:      8080,
//		Message:   "unexpected port opened",
//	})
package alert
