package watchdog_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watchdog"
)

func TestCheck_NoScansYet(t *testing.T) {
	w := watchdog.New(30*time.Second, "")
	s := w.Check()
	if s.Healthy {
		t.Fatal("expected unhealthy before any scan")
	}
	if s.Message == "" {
		t.Fatal("expected a message")
	}
}

func TestCheck_HealthyAfterScan(t *testing.T) {
	w := watchdog.New(30*time.Second, "")
	w.RecordScan()
	s := w.Check()
	if !s.Healthy {
		t.Fatalf("expected healthy, got: %s", s.Message)
	}
	if s.ScanCount != 1 {
		t.Fatalf("expected scan count 1, got %d", s.ScanCount)
	}
}

func TestCheck_StaleAfterThreshold(t *testing.T) {
	w := watchdog.New(1*time.Millisecond, "")
	w.RecordScan()
	time.Sleep(5 * time.Millisecond)
	s := w.Check()
	if s.Healthy {
		t.Fatal("expected unhealthy after staleness threshold")
	}
}

func TestRecordScan_IncrementsCount(t *testing.T) {
	w := watchdog.New(time.Minute, "")
	for i := 0; i < 5; i++ {
		w.RecordScan()
	}
	if w.Check().ScanCount != 5 {
		t.Fatalf("expected scan count 5")
	}
}

func TestWritePID_And_RemovePID(t *testing.T) {
	dir := t.TempDir()
	pidFile := filepath.Join(dir, "portwatch.pid")
	w := watchdog.New(time.Minute, pidFile)

	if err := w.WritePID(); err != nil {
		t.Fatalf("WritePID: %v", err)
	}
	data, err := os.ReadFile(pidFile)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty PID file")
	}

	w.RemovePID()
	if _, err := os.Stat(pidFile); !os.IsNotExist(err) {
		t.Fatal("expected PID file to be removed")
	}
}

func TestWritePID_EmptyPath(t *testing.T) {
	w := watchdog.New(time.Minute, "")
	if err := w.WritePID(); err != nil {
		t.Fatalf("expected no error for empty pid path, got: %v", err)
	}
}
