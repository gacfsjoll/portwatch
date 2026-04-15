// Package state provides persistence for port scan snapshots,
// allowing portwatch to detect changes across restarts.
package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Snapshot holds the result of a port scan at a point in time.
type Snapshot struct {
	Timestamp time.Time `json:"timestamp"`
	OpenPorts []uint16  `json:"open_ports"`
}

// Store persists and retrieves port scan snapshots to/from disk.
type Store struct {
	path string
}

// NewStore creates a Store that reads/writes snapshots at the given file path.
// The parent directory is created if it does not exist.
func NewStore(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	return &Store{path: path}, nil
}

// Save writes the snapshot to disk, overwriting any previous snapshot.
func (s *Store) Save(snap Snapshot) error {
	f, err := os.CreateTemp(filepath.Dir(s.path), ".portwatch-snap-*")
	if err != nil {
		return err
	}
	tmpName := f.Name()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		f.Close()
		os.Remove(tmpName)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}
	return os.Rename(tmpName, s.path)
}

// Load reads the most recent snapshot from disk.
// If no snapshot file exists, it returns an empty Snapshot and no error.
func (s *Store) Load() (Snapshot, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return Snapshot{}, nil
	}
	if err != nil {
		return Snapshot{}, err
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}
