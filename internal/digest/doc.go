// Package digest provides helpers for producing stable, order-independent
// SHA-256 digests of port-state snapshots.
//
// # Overview
//
// Digests allow portwatch to cheaply decide whether a freshly scanned snapshot
// is identical to the last persisted baseline before doing a more expensive
// field-by-field comparison via the snapshot package.
//
// Internally, each port entry is serialised to a canonical byte representation,
// the resulting byte slices are sorted, and a single SHA-256 hash is computed
// over the concatenation. This guarantees that the digest is independent of the
// order in which ports are reported by the underlying scanner.
//
// # Usage
//
//	sum, err := digest.Compute(ports)
//
//	changed, err := digest.Equal(previous, current)
//	if err != nil { ... }
//	if !changed {
//	    // nothing to do
//	}
package digest
