package metrics_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

var fixedTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func fixedClock() time.Time { return fixedTime }

func TestRecordScan_Increments(t *testing.T) {
	r := metrics.NewWithClock(fixedClock)
	r.RecordScan()
	r.RecordScan()
	s := r.Snapshot()
	if s.ScansTotal != 2 {
		t.Fatalf("expected 2 scans, got %d", s.ScansTotal)
	}
	if !s.LastScanTime.Equal(fixedTime) {
		t.Fatalf("unexpected last scan time: %v", s.LastScanTime)
	}
}

func TestRecordAlert_Increments(t *testing.T) {
	r := metrics.NewWithClock(fixedClock)
	r.RecordAlert(2, 1)
	r.RecordAlert(0, 3)
	s := r.Snapshot()
	if s.AlertsTotal != 2 {
		t.Fatalf("expected 2 alerts, got %d", s.AlertsTotal)
	}
	if s.OpenedTotal != 2 {
		t.Fatalf("expected 2 opened, got %d", s.OpenedTotal)
	}
	if s.ClosedTotal != 4 {
		t.Fatalf("expected 4 closed, got %d", s.ClosedTotal)
	}
}

func TestSnapshot_IsCopy(t *testing.T) {
	r := metrics.NewWithClock(fixedClock)
	r.RecordScan()
	s1 := r.Snapshot()
	r.RecordScan()
	s2 := r.Snapshot()
	if s1.ScansTotal == s2.ScansTotal {
		t.Fatal("snapshot should be independent copy")
	}
}

func TestPrint_ContainsFields(t *testing.T) {
	r := metrics.NewWithClock(fixedClock)
	r.RecordScan()
	r.RecordAlert(1, 0)
	var buf bytes.Buffer
	r.Print(&buf)
	out := buf.String()
	for _, want := range []string{"scans=1", "alerts=1", "opened=1", "closed=0"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q: %s", want, out)
		}
	}
}

func TestPrint_DefaultsToStdout(t *testing.T) {
	// Should not panic when w is nil.
	r := metrics.New()
	r.Print(nil)
}
