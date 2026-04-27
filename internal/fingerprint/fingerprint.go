// Package fingerprint produces a deterministic identity string for a port
// scan result set, suitable for change detection and deduplication.
package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/user/portwatch/internal/scanner"
)

// Fingerprint is a hex-encoded SHA-256 digest of a sorted port list.
type Fingerprint string

// String returns the underlying hex string.
func (f Fingerprint) String() string { return string(f) }

// Equal reports whether two fingerprints are identical.
func (f Fingerprint) Equal(other Fingerprint) bool { return f == other }

// Compute derives a Fingerprint from a slice of PortState values.
// The result is stable regardless of the order ports are supplied.
func Compute(ports []scanner.PortState) Fingerprint {
	nums := make([]int, len(ports))
	for i, p := range ports {
		nums[i] = p.Port
	}
	sort.Ints(nums)

	h := sha256.New()
	for _, n := range nums {
		fmt.Fprintf(h, "%d\n", n)
	}
	return Fingerprint(hex.EncodeToString(h.Sum(nil)))
}

// Changed returns true when before and after produce different fingerprints.
func Changed(before, after []scanner.PortState) bool {
	return !Compute(before).Equal(Compute(after))
}
