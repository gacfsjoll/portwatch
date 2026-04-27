package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func writeBackup(t *testing.T, dir, base, suffix string, age time.Duration) string {
	t.Helper()
	name := filepath.Join(dir, base+"."+suffix)
	if err := os.WriteFile(name, []byte("log data\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	mod := time.Now().Add(-age)
	if err := os.Chtimes(name, mod, mod); err != nil {
		t.Fatal(err)
	}
	return name
}

func TestRunRotatorInfo_NoBackups(t *testing.T) {
	dir := t.TempDir()
	logFile := filepath.Join(dir, "portwatch.log")
	_ = os.WriteFile(logFile, []byte("current\n"), 0o644)

	out := captureStdout(t, func() {
		if err := runRotatorInfo([]string{"--log", logFile}); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "no rotated backups") {
		t.Fatalf("expected no-backups message, got: %q", out)
	}
}

func TestRunRotatorInfo_ShowsBackups(t *testing.T) {
	dir := t.TempDir()
	logFile := filepath.Join(dir, "portwatch.log")
	writeBackup(t, dir, "portwatch.log", "20240601T120000Z", 2*time.Hour)

	out := captureStdout(t, func() {
		if err := runRotatorInfo([]string{"--log", logFile}); err != nil {
			t.Fatal(err)
		}
	})
	if !strings.Contains(out, "20240601T120000Z") {
		t.Fatalf("expected backup name in output, got: %q", out)
	}
}

func TestRunRotatorInfo_MissingLogFlag(t *testing.T) {
	err := runRotatorInfo([]string{})
	if err == nil || !strings.Contains(err.Error(), "--log is required") {
		t.Fatalf("expected --log required error, got: %v", err)
	}
}

func TestRunRotatorPrune_RemovesOldBackups(t *testing.T) {
	dir := t.TempDir()
	logFile := filepath.Join(dir, "portwatch.log")
	old := writeBackup(t, dir, "portwatch.log", "20230101T000000Z", 400*24*time.Hour)
	recent := writeBackup(t, dir, "portwatch.log", "20240601T000000Z", 1*time.Hour)

	_ = captureStdout(t, func() {
		if err := runRotatorPrune([]string{"--log", logFile, "--older-than", "168h"}); err != nil {
			t.Fatal(err)
		}
	})

	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Fatal("old backup should have been removed")
	}
	if _, err := os.Stat(recent); err != nil {
		t.Fatal("recent backup should still exist")
	}
}

func TestRunRotatorPrune_MissingLogFlag(t *testing.T) {
	err := runRotatorPrune([]string{})
	if err == nil || !strings.Contains(err.Error(), "--log is required") {
		t.Fatalf("expected --log required error, got: %v", err)
	}
}

// captureStdout redirects os.Stdout for the duration of fn and returns the
// captured output.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w
	fn()
	_ = w.Close()
	os.Stdout = old
	var sb strings.Builder
	buf := make([]byte, 4096)
	for {
		n, err := r.Read(buf)
		sb.Write(buf[:n])
		if err != nil {
			break
		}
	}
	return sb.String()
}
