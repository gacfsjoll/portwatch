// Package rotator implements size-based log file rotation for portwatch.
//
// A Rotator wraps a file path and satisfies io.WriteCloser. When a write
// would cause the file to exceed the configured MaxBytes threshold the
// current file is closed, renamed with a UTC timestamp suffix
// (e.g. portwatch.log.20240601T120000Z), and a fresh file is opened at the
// original path. Old backup files beyond MaxBackups are pruned automatically.
//
// Typical usage:
//
//	r, err := rotator.New("/var/log/portwatch/events.log", 10*1024*1024, 5)
//	if err != nil { ... }
//	defer r.Close()
//	log.SetOutput(r)
//
// Rotator is safe for concurrent use.
package rotator
