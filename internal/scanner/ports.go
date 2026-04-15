package scanner

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

// PortState represents the state of a single open port.
type PortState struct {
	Port     int
	Protocol string
	Address  string
}

// String returns a human-readable representation of the port state.
func (p PortState) String() string {
	return fmt.Sprintf("%s:%d/%s", p.Address, p.Port, p.Protocol)
}

// Scanner scans for open ports on the local machine.
type Scanner struct {
	StartPort int
	EndPort   int
	Timeout   time.Duration
}

// NewScanner creates a Scanner with sensible defaults.
func NewScanner(startPort, endPort int, timeout time.Duration) *Scanner {
	return &Scanner{
		StartPort: startPort,
		EndPort:   endPort,
		Timeout:   timeout,
	}
}

// Scan checks all ports in the configured range and returns those that are open.
func (s *Scanner) Scan() ([]PortState, error) {
	if s.StartPort < 1 || s.EndPort > 65535 || s.StartPort > s.EndPort {
		return nil, fmt.Errorf("invalid port range: %d-%d", s.StartPort, s.EndPort)
	}

	var open []PortState

	for port := s.StartPort; port <= s.EndPort; port++ {
		addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
		conn, err := net.DialTimeout("tcp", addr, s.Timeout)
		if err == nil {
			conn.Close()
			open = append(open, PortState{
				Port:     port,
				Protocol: "tcp",
				Address:  "127.0.0.1",
			})
		}
	}

	return open, nil
}
