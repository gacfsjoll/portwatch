package main

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/audit"
)

type fakeAuditCfg struct{ path string }

func (f *fakeAuditCfg) AuditLogPath() string { return f.path }

func TestRunAuditLog_Empty(t *testing.T) {
	cfg := &fakeAuditCfg{path: filepath.Join(t.TempDir(), "audit.jsonl")}
	if err := runAuditLog(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunAuditLog_WithEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.jsonl")
	l := audit.New(path)
	l.Log("cli", "acknowledge", 8080, "test")
	l.Log("daemon", "alert", 9090, "opened")

	cfg := &fakeAuditCfg{path: path}
	if err := runAuditLog(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunAuditLog_MissingFile(t *testing.T) {
	cfg := &fakeAuditCfg{path: "/nonexistent/audit.jsonl"}
	if err := runAuditLog(cfg); err != nil {
		t.Fatalf("missing file should not error: %v", err)
	}
}

func TestAuditEntry_ZeroPort(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.jsonl")
	l := audit.New(path)
	_ = l.Log("daemon", "scan", 0, "routine")

	entries, err := audit.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Port != 0 {
		t.Errorf("expected port 0, got %d", entries[0].Port)
	}
	_ = time.Now()
}
