package retention_test

import (
	"testing"
	"time"

	"github.com/example/portwatch/internal/retention"
)

var epoch = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func fixedNow() time.Time { return epoch }

func TestDefaultPolicy(t *testing.T) {
	p := retention.DefaultPolicy()
	if p.MaxAge != 30*24*time.Hour {
		t.Fatalf("expected 30d MaxAge, got %v", p.MaxAge)
	}
	if p.MaxEntries != 10_000 {
		t.Fatalf("expected 10000 MaxEntries, got %d", p.MaxEntries)
	}
}

func TestApply_FiltersOldEntries(t *testing.T) {
	pr := retention.NewWithClock(retention.Policy{MaxAge: 24 * time.Hour, MaxEntries: 0}, fixedNow)

	entries := []time.Time{
		epoch.Add(-48 * time.Hour), // too old
		epoch.Add(-25 * time.Hour), // too old
		epoch.Add(-23 * time.Hour), // keep
		epoch.Add(-1 * time.Hour),  // keep
	}

	got := pr.Apply(entries)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestApply_RespectsMaxEntries(t *testing.T) {
	pr := retention.NewWithClock(retention.Policy{MaxAge: 24 * time.Hour, MaxEntries: 2}, fixedNow)

	entries := []time.Time{
		epoch.Add(-20 * time.Hour),
		epoch.Add(-10 * time.Hour),
		epoch.Add(-5 * time.Hour),
		epoch.Add(-1 * time.Hour),
	}

	got := pr.Apply(entries)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries after MaxEntries cap, got %d", len(got))
	}
	// Should keep the two most recent.
	if !got[0].Equal(epoch.Add(-5*time.Hour)) || !got[1].Equal(epoch.Add(-1*time.Hour)) {
		t.Fatalf("unexpected entries retained: %v", got)
	}
}

func TestApply_EmptyInput(t *testing.T) {
	pr := retention.NewWithClock(retention.DefaultPolicy(), fixedNow)
	got := pr.Apply(nil)
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %d entries", len(got))
	}
}

func TestShouldPrune(t *testing.T) {
	pr := retention.NewWithClock(retention.Policy{MaxAge: 24 * time.Hour}, fixedNow)

	if !pr.ShouldPrune(epoch.Add(-25 * time.Hour)) {
		t.Error("expected old entry to be prunable")
	}
	if pr.ShouldPrune(epoch.Add(-1 * time.Hour)) {
		t.Error("expected recent entry to be kept")
	}
}
