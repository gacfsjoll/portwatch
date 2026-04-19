// Package audit provides a structured audit log for portwatch events,
// recording who or what triggered an action and when.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Actor     string    `json:"actor"`
	Action    string    `json:"action"`
	Port      int       `json:"port,omitempty"`
	Detail    string    `json:"detail,omitempty"`
}

// Logger writes audit entries to a newline-delimited JSON file.
type Logger struct {
	mu   sync.Mutex
	path string
	now  func() time.Time
}

// New returns a Logger that appends to the file at path.
func New(path string) *Logger {
	return &Logger{path: path, now: time.Now}
}

// Log appends an entry to the audit log.
func (l *Logger) Log(actor, action string, port int, detail string) error {
	e := Entry{
		Timestamp: l.now().UTC(),
		Actor:     actor,
		Action:    action,
		Port:      port,
		Detail:    detail,
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal: %w", err)
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if err := os.MkdirAll(filepath.Dir(l.path), 0o755); err != nil {
		return fmt.Errorf("audit: mkdir: %w", err)
	}
	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("audit: open: %w", err)
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%s\n", b)
	return err
}

// Load reads all audit entries from path.
func Load(path string) ([]Entry, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("audit: read: %w", err)
	}
	var entries []Entry
	for _, line := range splitLines(data) {
		if len(line) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, fmt.Errorf("audit: unmarshal: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
