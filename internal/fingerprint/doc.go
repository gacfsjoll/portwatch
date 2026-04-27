// Package fingerprint derives a stable, deterministic identity for a set of
// open ports observed during a single scan cycle.
//
// # Overview
//
// A Fingerprint is a hex-encoded SHA-256 digest computed over the sorted list
// of open port numbers. Because ports are sorted before hashing, the result is
// independent of the order in which the scanner returns them.
//
// # Usage
//
//	f := fingerprint.Compute(ports)
//	if fingerprint.Changed(previous, current) {
//		// trigger alert
//	}
//
// Fingerprints are cheap to store alongside snapshots and history entries,
// enabling fast equality checks without re-diffing full port lists.
package fingerprint
