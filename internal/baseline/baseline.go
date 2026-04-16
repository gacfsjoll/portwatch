// Package baseline manages the trusted set of open ports that portwatch
// considers "expected". Any deviation from the baseline triggers an alert.
package baseline

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Baseline holds the set of ports that are considered normal.
type Baseline struct {
	Ports     []int     `json:"ports"`
	CapturedAt time.Time `json:"captured_at"`
}

// Contains reports whether port p is part of the baseline.
func (b *Baseline) Contains(p int) bool {
	for _, bp := range b.Ports {
		if bp == p {
			return true
		}
	}
	return false
}

// Manager handles persistence and retrieval of the port baseline.
type Manager struct {
	path string
}

// NewManager returns a Manager that stores the baseline at the given path.
func NewManager(path string) *Manager {
	return &Manager{path: path}
}

// Save writes the provided port list as the new baseline to disk.
func (m *Manager) Save(ports []int) error {
	sorted := make([]int, len(ports))
	copy(sorted, ports)
	sort.Ints(sorted)

	b := Baseline{
		Ports:     sorted,
		CapturedAt: time.Now().UTC(),
	}

	if err := os.MkdirAll(filepath.Dir(m.path), 0o755); err != nil {
		return err
	}

	f, err := os.Create(m.path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(b)
}

// Load reads the baseline from disk. Returns ErrNoBaseline if the file does
// not exist yet.
func (m *Manager) Load() (*Baseline, error) {
	f, err := os.Open(m.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNoBaseline
		}
		return nil, err
	}
	defer f.Close()

	var b Baseline
	if err := json.NewDecoder(f).Decode(&b); err != nil {
		return nil, err
	}
	return &b, nil
}

// ErrNoBaseline is returned when no baseline file has been captured yet.
var ErrNoBaseline = errors.New("no baseline captured; run 'portwatch capture' first")
