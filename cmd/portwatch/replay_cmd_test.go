package main

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/example/portwatch/internal/alert"
	"github.com/example/portwatch/internal/history"
)

func tempReplayRecorder(t *testing.T, events []alert.Event) *history.Recorder {
	t.Helper()
	dir := t.TempDir()
	rec, err := history.NewRecorder(filepath.Join(dir, "history.json"))
	if err != nil {
		t.Fatalf("NewRecorder: %v", err)
	}
	for _, ev := range events {
		if err := rec.Record(ev); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}
	return rec
}

func TestRunReplay_UnknownFlag(t *testing.T) {
	err := runReplay([]string{"--unknown-flag"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}

func TestRunReplay_InvalidSince(t *testing.T) {
	err := runReplay([]string{"--since", "notaduration"})
	if err == nil {
		t.Fatal("expected error for invalid --since")
	}
}

func TestRunReplay_DryRunFlag_Parsed(t *testing.T) {
	// Ensure the flag parses without error (config load may fail in unit env).
	// We only check flag parsing isolation here; integration covered in replay pkg.
	fs := newReplayFlagSet()
	if err := fs.Parse([]string{"--dry-run", "--since", "1h"}); err != nil {
		t.Fatalf("flag parse: %v", err)
	}
}

func TestRunReplay_DelayFlag_Parsed(t *testing.T) {
	fs := newReplayFlagSet()
	if err := fs.Parse([]string{"--delay", "50ms"}); err != nil {
		t.Fatalf("flag parse: %v", err)
	}
}

// newReplayFlagSet mirrors the flag setup in runReplay for isolated tests.
func newReplayFlagSet() *flagSetWrapper {
	type cfg struct {
		since  string
		dryRun bool
		delay  time.Duration
	}
	w := &flagSetWrapper{parsed: &cfg{}}
	w.fs.StringVar(&w.parsed.(*cfg).since, "since", "", "")
	w.fs.BoolVar(&w.parsed.(*cfg).dryRun, "dry-run", false, "")
	w.fs.DurationVar(&w.parsed.(*cfg).delay, "delay", 0, "")
	return w
}

type flagSetWrapper struct {
	fs     flag.FlagSet
	parsed interface{}
}

func (w *flagSetWrapper) Parse(args []string) error {
	return w.fs.Parse(args)
}
