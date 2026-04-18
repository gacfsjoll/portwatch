// Package acknowledge tracks ports that have been acknowledged by the operator,
// suppressing repeated alerts until the port state changes again.
package acknowledge

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// Store persists acknowledged ports to disk.
type Store struct {
	mu   sync.Mutex
	path string
	acks map[uint16]struct{}
}

// NewStore creates a new Store backed by the given file path.
func NewStore(path string) *Store {
	return &Store{path: path, acks: make(map[uint16]struct{})}
}

// Load reads acknowledged ports from disk. Missing file is not an error.
func (s *Store) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	var ports []uint16
	if err := json.Unmarshal(data, &ports); err != nil {
		return err
	}
	s.acks = make(map[uint16]struct{}, len(ports))
	for _, p := range ports {
		s.acks[p] = struct{}{}
	}
	return nil
}

// Acknowledge marks a port as acknowledged.
func (s *Store) Acknowledge(port uint16) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.acks[port] = struct{}{}
	return s.save()
}

// Revoke removes the acknowledgement for a port.
func (s *Store) Revoke(port uint16) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.acks, port)
	return s.save()
}

// IsAcknowledged returns true if the port has been acknowledged.
func (s *Store) IsAcknowledged(port uint16) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.acks[port]
	return ok
}

// List returns all currently acknowledged ports.
func (s *Store) List() []uint16 {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]uint16, 0, len(s.acks))
	for p := range s.acks {
		out = append(out, p)
	}
	return out
}

func (s *Store) save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	ports := make([]uint16, 0, len(s.acks))
	for p := range s.acks {
		ports = append(ports, p)
	}
	data, err := json.MarshalIndent(ports, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
