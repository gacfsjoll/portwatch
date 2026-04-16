package history_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
)

func tempHistoryPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "sub", "history.jsonl")
}

func TestRecorder_RecordAndLoad(t *testing.T) {
	path := tempHistoryPath(t)
	rec, err := history.NewRecorder(path)
	if err != nil {
		t.Fatalf("NewRecorder: %v", err)
	}

	entries := []history.Entry{
		{Timestamp: time.Now().UTC(), Port: 8080, Proto: "tcp", Event: "opened"},
		{Timestamp: time.Now().UTC(), Port: 443, Proto: "tcp", Event: "closed"},
	}
	for _, e := range entries {
		if err := rec.Record(e); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}

	loaded, err := history.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(loaded))
	}
	for i, e := range entries {
		if loaded[i].Port != e.Port || loaded[i].Event != e.Event {
			t.Errorf("entry %d mismatch: got %+v, want %+v", i, loaded[i], e)
		}
	}
}

func TestLoad_MissingFile(t *testing.T) {
	entries, err := history.Load("/nonexistent/path/history.jsonl")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if entries != nil {
		t.Fatalf("expected nil entries, got %v", entries)
	}
}

func TestRecorder_CreatesParentDirs(t *testing.T) {
	path := tempHistoryPath(t)
	_, err := history.NewRecorder(path)
	if err != nil {
		t.Fatalf("NewRecorder: %v", err)
	}
	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		t.Fatalf("expected parent dirs to exist: %v", err)
	}
}
