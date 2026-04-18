// Package digest provides helpers for producing stable, order-independent
// SHA-256 digests of port-state snapshots.
//
// # Overview
//
// Digests allow portwatch to cheaply decide whether a freshly scanned snapshot
// is identical to the last persisted baseline before doing a more expensive
// field-by-field comparison via the snapshot package.
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
