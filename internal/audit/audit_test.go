package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempAuditPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "audit", "audit.jsonl")
}

func fixedNow() func() time.Time {
	ts := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	return func() time.Time { return ts }
}

func TestLog_And_Load(t *testing.T) {
	path := tempAuditPath(t)
	l := New(path)
	l.now = fixedNow()

	if err := l.Log("cli", "acknowledge", 8080, "user added"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := l.Log("daemon", "alert", 9090, "port opened"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := Load(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Action != "acknowledge" {
		t.Errorf("expected acknowledge, got %s", entries[0].Action)
	}
	if entries[1].Port != 9090 {
		t.Errorf("expected port 9090, got %d", entries[1].Port)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	entries, err := Load("/nonexistent/path/audit.jsonl")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries, got %v", entries)
	}
}

func TestLog_CreatesParentDirs(t *testing.T) {
	path := tempAuditPath(t)
	l := New(path)
	if err := l.Log("cli", "suppress", 443, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}

func TestEntry_Timestamp(t *testing.T) {
	path := tempAuditPath(t)
	l := New(path)
	l.now = fixedNow()
	_ = l.Log("daemon", "scan", 0, "")

	entries, _ := Load(path)
	expected := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	if !entries[0].Timestamp.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, entries[0].Timestamp)
	}
}
