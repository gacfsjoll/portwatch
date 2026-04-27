package rotator_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/rotator"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "portwatch.log")
}

func TestNew_CreatesFile(t *testing.T) {
	p := tempPath(t)
	r, err := rotator.New(p, 1024, 3)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func TestWrite_AppendsData(t *testing.T) {
	p := tempPath(t)
	r, _ := rotator.New(p, 1024, 3)
	defer r.Close()

	_, err := fmt.Fprintln(r, "hello portwatch")
	if err != nil {
		t.Fatalf("Write: %v", err)
	}
	data, _ := os.ReadFile(p)
	if string(data) != "hello portwatch\n" {
		t.Fatalf("unexpected content: %q", data)
	}
}

func TestWrite_RotatesWhenFull(t *testing.T) {
	p := tempPath(t)
	r, _ := rotator.New(p, 20, 5)
	defer r.Close()

	// Write enough to trigger rotation.
	for i := 0; i < 5; i++ {
		_, err := fmt.Fprintln(r, "0123456789") // 11 bytes each
		if err != nil {
			t.Fatalf("Write[%d]: %v", i, err)
		}
	}

	matches, _ := filepath.Glob(p + ".*")
	if len(matches) == 0 {
		t.Fatal("expected at least one rotated backup, got none")
	}
}

func TestPruneBackups_RemovesOldest(t *testing.T) {
	p := tempPath(t)
	// maxBackups = 2 so only 2 backups are kept.
	r, _ := rotator.New(p, 10, 2)
	defer r.Close()

	for i := 0; i < 10; i++ {
		_, _ = fmt.Fprintln(r, "0123456789")
		time.Sleep(2 * time.Millisecond) // ensure distinct timestamps
	}

	matches, _ := filepath.Glob(p + ".*")
	if len(matches) > 2 {
		t.Fatalf("expected ≤2 backups, got %d", len(matches))
	}
}

func TestNew_DefaultsOnZeroValues(t *testing.T) {
	p := tempPath(t)
	r, err := rotator.New(p, 0, 0)
	if err != nil {
		t.Fatalf("New with zero values: %v", err)
	}
	r.Close()
}

func TestNew_CreatesParentDirs(t *testing.T) {
	base := t.TempDir()
	p := filepath.Join(base, "deep", "nested", "portwatch.log")
	r, err := rotator.New(p, 1024, 3)
	if err != nil {
		t.Fatalf("New nested path: %v", err)
	}
	defer r.Close()
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("file not created in nested dir: %v", err)
	}
}
