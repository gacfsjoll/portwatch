// Package digest computes and compares port-state digests so portwatch can
// detect whether the current snapshot differs from a previously persisted one
// without storing the full snapshot twice.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/user/portwatch/internal/scanner"
)

// Compute returns a stable hex-encoded SHA-256 digest of the given port list.
// The list is sorted before hashing so order does not affect the result.
func Compute(ports []scanner.PortState) (string, error) {
	sorted := make([]scanner.PortState, len(ports))
	copy(sorted, ports)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Port != sorted[j].Port {
			return sorted[i].Port < sorted[j].Port
		}
		return sorted[i].Proto < sorted[j].Proto
	})

	b, err := json.Marshal(sorted)
	if err != nil {
		return "", fmt.Errorf("digest: marshal: %w", err)
	}

	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}

// Equal returns true when the two port lists produce the same digest.
// An error is returned if either list cannot be hashed.
func Equal(a, b []scanner.PortState) (bool, error) {
	da, err := Compute(a)
	if err != nil {
		return false, err
	}
	db, err := Compute(b)
	if err != nil {
		return false, err
	}
	return da == db, nil
}
