// Package retention implements a configurable data-retention policy for
// portwatch's on-disk state.
//
// As portwatch runs continuously it accumulates history entries, audit log
// lines, and other timestamped records. Without pruning, these files grow
// without bound. The retention package provides a [Policy] type that
// captures the desired MaxAge and MaxEntries constraints, and a [Pruner]
// that applies those constraints to any slice of [time.Time] values.
//
// # Usage
//
//	// Build a pruner from the default policy.
//	pr := retention.New(retention.DefaultPolicy())
//
//	// Filter a slice of event timestamps.
//	keep := pr.Apply(timestamps)
//
// The [Pruner.ShouldPrune] helper is available for single-record checks
// when iterating over a data store that does not support bulk replacement.
package retention
