package state_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/state"
)

func tempStorePath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "sub", "state.json")
}

func TestStore_SaveAndLoad(t *testing.T) {
	path := tempStorePath(t)
	store, err := state.NewStore(path)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	snap := state.Snapshot{
		Timestamp: time.Now().UTC().Truncate(time.Second),
		OpenPorts: []uint16{22, 80, 443},
	}
	if err := store.Save(snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !got.Timestamp.Equal(snap.Timestamp) {
		t.Errorf("Timestamp: got %v, want %v", got.Timestamp, snap.Timestamp)
	}
	if len(got.OpenPorts) != len(snap.OpenPorts) {
		t.Fatalf("OpenPorts len: got %d, want %d", len(got.OpenPorts), len(snap.OpenPorts))
	}
	for i, p := range snap.OpenPorts {
		if got.OpenPorts[i] != p {
			t.Errorf("OpenPorts[%d]: got %d, want %d", i, got.OpenPorts[i], p)
		}
	}
}

func TestStore_LoadMissingFile(t *testing.T) {
	path := tempStorePath(t)
	store, err := state.NewStore(path)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	snap, err := store.Load()
	if err != nil {
		t.Fatalf("Load on missing file should not error: %v", err)
	}
	if len(snap.OpenPorts) != 0 {
		t.Errorf("expected empty snapshot, got %v", snap.OpenPorts)
	}
}

func TestStore_SaveCreatesParentDirs(t *testing.T) {
	path := tempStorePath(t)
	store, err := state.NewStore(path)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	if err := store.Save(state.Snapshot{OpenPorts: []uint16{8080}}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected state file to exist: %v", err)
	}
}
