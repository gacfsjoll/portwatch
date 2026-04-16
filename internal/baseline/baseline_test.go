package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/baseline"
)

func tempBaselinePath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "sub", "baseline.json")
}

func TestManager_SaveAndLoad(t *testing.T) {
	path := tempBaselinePath(t)
	m := baseline.NewManager(path)

	ports := []int{8080, 443, 80, 22}
	if err := m.Save(ports); err != nil {
		t.Fatalf("Save: %v", err)
	}

	b, err := m.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	// Expect ports to be sorted.
	want := []int{22, 80, 443, 8080}
	if len(b.Ports) != len(want) {
		t.Fatalf("got %d ports, want %d", len(b.Ports), len(want))
	}
	for i, p := range want {
		if b.Ports[i] != p {
			t.Errorf("ports[%d] = %d, want %d", i, b.Ports[i], p)
		}
	}

	if b.CapturedAt.IsZero() {
		t.Error("CapturedAt should not be zero")
	}
}

func TestManager_LoadMissingFile(t *testing.T) {
	m := baseline.NewManager("/nonexistent/path/baseline.json")
	_, err := m.Load()
	if err != baseline.ErrNoBaseline {
		t.Errorf("expected ErrNoBaseline, got %v", err)
	}
}

func TestManager_SaveCreatesParentDirs(t *testing.T) {
	path := tempBaselinePath(t)
	m := baseline.NewManager(path)

	if err := m.Save([]int{9000}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}

func TestBaseline_Contains(t *testing.T) {
	path := tempBaselinePath(t)
	m := baseline.NewManager(path)

	_ = m.Save([]int{22, 80, 443})
	b, _ := m.Load()

	if !b.Contains(80) {
		t.Error("expected baseline to contain port 80")
	}
	if b.Contains(8080) {
		t.Error("expected baseline NOT to contain port 8080")
	}
}
