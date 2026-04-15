package scanner

import (
	"net"
	"testing"
	"time"
)

// startTestListener opens a TCP listener on an OS-assigned port and returns
// the port number along with a cleanup function.
func startTestListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return port, func() { ln.Close() }
}

func TestScan_DetectsOpenPort(t *testing.T) {
	port, cleanup := startTestListener(t)
	defer cleanup()

	s := NewScanner(port, port, 200*time.Millisecond)
	results, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 open port, got %d", len(results))
	}
	if results[0].Port != port {
		t.Errorf("expected port %d, got %d", port, results[0].Port)
	}
	if results[0].Protocol != "tcp" {
		t.Errorf("expected protocol tcp, got %s", results[0].Protocol)
	}
}

func TestScan_InvalidRange(t *testing.T) {
	s := NewScanner(9000, 8000, 200*time.Millisecond)
	_, err := s.Scan()
	if err == nil {
		t.Fatal("expected error for invalid port range, got nil")
	}
}

func TestPortState_String(t *testing.T) {
	p := PortState{Port: 8080, Protocol: "tcp", Address: "127.0.0.1"}
	expected := "127.0.0.1:8080/tcp"
	if p.String() != expected {
		t.Errorf("expected %q, got %q", expected, p.String())
	}
}
