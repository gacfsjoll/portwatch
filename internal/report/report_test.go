package report_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/report"
)

func makeEvents() []history.Event {
	return []history.Event{
		{Port: 80, Proto: "tcp", Kind: "opened", At: time.Now()},
		{Port: 443, Proto: "tcp", Kind: "opened", At: time.Now()},
		{Port: 8080, Proto: "tcp", Kind: "closed", At: time.Now()},
	}
}

func TestFromHistory_Counts(t *testing.T) {
	s := report.FromHistory(makeEvents())
	if s.TotalEvents != 3 {
		t.Fatalf("expected 3 events, got %d", s.TotalEvents)
	}
	if len(s.Opened) != 2 {
		t.Fatalf("expected 2 opened, got %d", len(s.Opened))
	}
	if len(s.Closed) != 1 {
		t.Fatalf("expected 1 closed, got %d", len(s.Closed))
	}
}

func TestFromHistory_Empty(t *testing.T) {
	s := report.FromHistory(nil)
	if s.TotalEvents != 0 {
		t.Fatalf("expected 0 events")
	}
}

func TestGenerator_Print(t *testing.T) {
	var buf bytes.Buffer
	g := report.NewGenerator(&buf)
	s := report.FromHistory(makeEvents())
	g.Print(s)
	out := buf.String()
	if !strings.Contains(out, "80/tcp") {
		t.Error("expected 80/tcp in output")
	}
	if !strings.Contains(out, "8080/tcp") {
		t.Error("expected 8080/tcp in output")
	}
	if !strings.Contains(out, "Total events") {
		t.Error("expected header in output")
	}
}

func TestGenerator_DefaultsToStdout(t *testing.T) {
	g := report.NewGenerator(nil)
	if g.Out == nil {
		t.Fatal("expected non-nil writer")
	}
}
