// Package history provides persistent recording and retrieval of port change
// events observed by portwatch.
//
// Events are stored as newline-delimited JSON (NDJSON) so the file can be
// appended to efficiently and streamed back without loading it entirely into
// memory.
//
// # Usage
//
//	rec, err := history.NewRecorder("/var/lib/portwatch/history.jsonl")
//	if err != nil { ... }
//
//	err = rec.Record(history.Entry{
//		Timestamp: time.Now().UTC(),
//		Port:      8080,
//		Proto:     "tcp",
//		Event:     "opened",
//	})
//
//	entries, err := history.Load("/var/lib/portwatch/history.jsonl")
package history
