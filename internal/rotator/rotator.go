// Package rotator provides log-file rotation for portwatch output files.
// It wraps an underlying file writer and rotates when the file exceeds a
// configured size threshold, keeping a configurable number of backups.
package rotator

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	DefaultMaxBytes   = 10 * 1024 * 1024 // 10 MiB
	DefaultMaxBackups = 5
)

// Rotator is a write-closer that rotates the underlying file when it grows
// beyond MaxBytes. Rotated files are renamed with a UTC timestamp suffix.
type Rotator struct {
	mu         sync.Mutex
	path       string
	maxBytes   int64
	maxBackups int
	file       *os.File
	size       int64
	now        func() time.Time
}

// New creates a Rotator for the given path.
func New(path string, maxBytes int64, maxBackups int) (*Rotator, error) {
	if maxBytes <= 0 {
		maxBytes = DefaultMaxBytes
	}
	if maxBackups <= 0 {
		maxBackups = DefaultMaxBackups
	}
	r := &Rotator{
		path:       path,
		maxBytes:   maxBytes,
		maxBackups: maxBackups,
		now:        time.Now,
	}
	if err := r.openOrCreate(); err != nil {
		return nil, err
	}
	return r, nil
}

// Write implements io.Writer. It rotates the file if the write would exceed
// MaxBytes.
func (r *Rotator) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.size+int64(len(p)) > r.maxBytes {
		if err := r.rotate(); err != nil {
			return 0, fmt.Errorf("rotator: rotate: %w", err)
		}
	}
	n, err := r.file.Write(p)
	r.size += int64(n)
	return n, err
}

// Close flushes and closes the underlying file.
func (r *Rotator) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.file.Close()
}

func (r *Rotator) openOrCreate() error {
	if err := os.MkdirAll(filepath.Dir(r.path), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(r.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	fi, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return err
	}
	r.file = f
	r.size = fi.Size()
	return nil
}

func (r *Rotator) rotate() error {
	if err := r.file.Close(); err != nil {
		return err
	}
	stamp := r.now().UTC().Format("20060102T150405Z")
	dst := fmt.Sprintf("%s.%s", r.path, stamp)
	if err := os.Rename(r.path, dst); err != nil {
		return err
	}
	r.pruneBackups()
	return r.openOrCreate()
}

func (r *Rotator) pruneBackups() {
	pattern := r.path + ".*"
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) <= r.maxBackups {
		return
	}
	// matches are lexicographically sorted; oldest have smallest timestamp.
	for _, old := range matches[:len(matches)-r.maxBackups] {
		_ = os.Remove(old)
	}
}
