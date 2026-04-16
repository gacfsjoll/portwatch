// Package history records port change events to a persistent JSON file.
package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Event represents a single port change occurrence.
type Event struct {
	Port  int       `json:"port"`
	Proto string    `json:"proto"`
	Kind  string    `json:"kind"` // "opened" | "closed"
	At    time.Time `json:"at"`
}

// Recorder appends events to a JSON file.
type Recorder struct {
	mu   sync.Mutex
	path string
}

// NewRecorder creates a Recorder that writes to path.
func NewRecorder(path string) *Recorder {
	return &Recorder{path: path}
}

// Record appends e to the history file.
func (r *Recorder) Record(e Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	events, err := load(r.path)
	if err != nil {
		return err
	}
	events = append(events, e)
	return save(r.path, events)
}

// Load reads all events from path. Returns empty slice if file is missing.
func Load(path string) ([]Event, error) {
	return load(path)
}

func load(path string) ([]Event, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []Event{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("history read: %w", err)
	}
	var events []Event
	if err := json.Unmarshal(data, &events); err != nil {
		return nil, fmt.Errorf("history parse: %w", err)
	}
	return events, nil
}

func save(path string, events []Event) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("history mkdir: %w", err)
	}
	data, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return fmt.Errorf("history marshal: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}
