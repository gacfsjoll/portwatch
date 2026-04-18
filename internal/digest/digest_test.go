package digest_test

import (
	"testing"

	"github.com/user/portwatch/internal/digest"
	"github.com/user/portwatch/internal/scanner"
)

func ps(port int, proto string) scanner.PortState {
	return scanner.PortState{Port: port, Proto: proto, Open: true}
}

func TestCompute_Stable(t *testing.T) {
	ports := []scanner.PortState{ps(443, "tcp"), ps(80, "tcp"), ps(22, "tcp")}
	d1, err := digest.Compute(ports)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// reverse order — digest must be identical
	reversed := []scanner.PortState{ps(22, "tcp"), ps(80, "tcp"), ps(443, "tcp")}
	d2, err := digest.Compute(reversed)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d1 != d2 {
		t.Errorf("expected equal digests, got %q vs %q", d1, d2)
	}
}

func TestCompute_DifferentPorts(t *testing.T) {
	d1, _ := digest.Compute([]scanner.PortState{ps(80, "tcp")})
	d2, _ := digest.Compute([]scanner.PortState{ps(8080, "tcp")})
	if d1 == d2 {
		t.Error("expected different digests for different ports")
	}
}

func TestCompute_Empty(t *testing.T) {
	d, err := digest.Compute(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == "" {
		t.Error("expected non-empty digest for nil input")
	}
}

func TestEqual_SamePorts(t *testing.T) {
	a := []scanner.PortState{ps(22, "tcp"), ps(80, "tcp")}
	b := []scanner.PortState{ps(80, "tcp"), ps(22, "tcp")}
	eq, err := digest.Equal(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !eq {
		t.Error("expected Equal to return true")
	}
}

func TestEqual_DifferentPorts(t *testing.T) {
	a := []scanner.PortState{ps(22, "tcp")}
	b := []scanner.PortState{ps(22, "tcp"), ps(443, "tcp")}
	eq, err := digest.Equal(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if eq {
		t.Error("expected Equal to return false")
	}
}
