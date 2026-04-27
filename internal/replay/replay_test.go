package replay_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/portwatch/internal/alert"
	"github.com/example/portwatch/internal/history"
	"github.com/example/portwatch/internal/replay"
)

func tempRecorder(t *testing.T) *history.Recorder {
	t.Helper()
	dir := t.TempDir()
	rec, err := history.NewRecorder(filepath.Join(dir, "history.json"))
	if err != nil {
		t.Fatalf("NewRecorder: %v", err)
	}
	return rec
}

type captureNotifier struct {
	events []alert.Event
	errOn  int // return error on this call (1-indexed), 0 = never
	calls  int
}

func (c *captureNotifier) Notify(ev alert.Event) error {
	c.calls++
	if c.errOn > 0 && c.calls == c.errOn {
		return errors.New("injected error")
	}
	c.events = append(c.events, ev)
	return nil
}

func seedEvents(t *testing.T, rec *history.Recorder, ports []int) {
	t.Helper()
	for _, p := range ports {
		if err := rec.Record(alert.Event{Port: p, Kind: "opened", Time: time.Now().UTC()}); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}
}

func TestReplay_DispatchesAllEvents(t *testing.T) {
	rec := tempRecorder(t)
	seedEvents(t, rec, []int{80, 443, 8080})

	n := &captureNotifier{}
	r := replay.New(n, replay.Options{})
	count, err := r.Run(context.Background(), rec)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if count != 3 {
		t.Errorf("count = %d, want 3", count)
	}
	if len(n.events) != 3 {
		t.Errorf("dispatched %d events, want 3", len(n.events))
	}
}

func TestReplay_DryRunDoesNotNotify(t *testing.T) {
	rec := tempRecorder(t)
	seedEvents(t, rec, []int{22, 80})

	n := &captureNotifier{}
	r := replay.New(n, replay.Options{DryRun: true})
	count, err := r.Run(context.Background(), rec)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
	if len(n.events) != 0 {
		t.Errorf("expected no dispatched events in dry-run, got %d", len(n.events))
	}
}

func TestReplay_SinceFiltersOldEvents(t *testing.T) {
	rec := tempRecorder(t)
	// Record an old event manually by manipulating the recorder via Record.
	// We rely on the Since window being very short.
	old := alert.Event{Port: 22, Kind: "opened", Time: time.Now().UTC().Add(-2 * time.Hour)}
	recent := alert.Event{Port: 80, Kind: "opened", Time: time.Now().UTC()}
	_ = rec.Record(old)
	_ = rec.Record(recent)

	n := &captureNotifier{}
	r := replay.New(n, replay.Options{Since: time.Hour})
	count, err := r.Run(context.Background(), rec)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if count != 1 {
		t.Errorf("count = %d, want 1", count)
	}
	if len(n.events) != 1 || n.events[0].Port != 80 {
		t.Errorf("expected port 80, got %+v", n.events)
	}
}

func TestReplay_NotifyError_Aborts(t *testing.T) {
	rec := tempRecorder(t)
	seedEvents(t, rec, []int{80, 443, 8080})

	n := &captureNotifier{errOn: 2}
	r := replay.New(n, replay.Options{})
	_, err := r.Run(context.Background(), rec)
	if err == nil {
		t.Fatal("expected error from notifier, got nil")
	}
}

func TestReplay_EmptyHistory(t *testing.T) {
	dir := t.TempDir()
	_ = os.Remove(filepath.Join(dir, "history.json"))
	rec, _ := history.NewRecorder(filepath.Join(dir, "history.json"))

	n := &captureNotifier{}
	r := replay.New(n, replay.Options{})
	count, err := r.Run(context.Background(), rec)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if count != 0 {
		t.Errorf("count = %d, want 0", count)
	}
}
