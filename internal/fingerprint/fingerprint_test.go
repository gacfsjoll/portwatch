package fingerprint_test

import (
	"testing"

	"github.com/user/portwatch/internal/fingerprint"
	"github.com/user/portwatch/internal/scanner"
)

func ps(ports ...int) []scanner.PortState {
	out := make([]scanner.PortState, len(ports))
	for i, p := range ports {
		out[i] = scanner.PortState{Port: p, Open: true}
	}
	return out
}

func TestCompute_Stable(t *testing.T) {
	a := fingerprint.Compute(ps(80, 443, 22))
	b := fingerprint.Compute(ps(443, 22, 80))
	if !a.Equal(b) {
		t.Fatalf("expected stable fingerprint, got %s vs %s", a, b)
	}
}

func TestCompute_Empty(t *testing.T) {
	f := fingerprint.Compute(nil)
	if f.String() == "" {
		t.Fatal("expected non-empty fingerprint for empty port list")
	}
}

func TestCompute_DifferentPorts(t *testing.T) {
	a := fingerprint.Compute(ps(80))
	b := fingerprint.Compute(ps(443))
	if a.Equal(b) {
		t.Fatal("different ports should produce different fingerprints")
	}
}

func TestChanged_DetectsChange(t *testing.T) {
	if !fingerprint.Changed(ps(80), ps(80, 443)) {
		t.Fatal("expected Changed to return true when ports differ")
	}
}

func TestChanged_NoChange(t *testing.T) {
	if fingerprint.Changed(ps(22, 80), ps(80, 22)) {
		t.Fatal("expected Changed to return false for same ports in different order")
	}
}

func TestFingerprint_String(t *testing.T) {
	f := fingerprint.Compute(ps(8080))
	if len(f.String()) != 64 {
		t.Fatalf("expected 64-char hex string, got len %d", len(f.String()))
	}
}
