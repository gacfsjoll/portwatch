// Package history records port change events to a persistent log.
package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Entry represents a single recorded port change event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Port      uint16    `json:"port"`
	Proto     string    `json:"proto"`
	Event     string    `json:"event"` // "opened" or "closed"
}

// Recorder appends history entries to a newline-delimited JSON file.
type Recorder struct {
	mu   sync.Mutex
	path string
}

// NewRecorder creates a Recorder that writes to the given file path.
func NewRecorder(path string) (*Recorder, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("history: create dirs: %w", err)
	}
	return &Recorder{path: path}, nil
}

// Record appends an entry to the history file.
func (r *Recorder) Record(e Entry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	f, err := os.OpenFile(r.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("history: open file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(e); err != nil {
		return fmt.Errorf("history: encode entry: %w", err)
	}
	return nil
}

// Load reads all entries from the history file.
func Load(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("history: open file: %w", err)
	}
	defer f.Close()

	var entries []Entry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("history: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
