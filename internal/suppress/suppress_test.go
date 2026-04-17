package suppress

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time { return func() time.Time { return t } }

func TestIsSuppressed_ActiveEntry(t *testing.T) {
	base := time.Now()
	l := New()
	l.now = fixedNow(base)
	l.Suppress(8080, 5*time.Minute, "maintenance")
	if !l.IsSuppressed(8080) {
		t.Fatal("expected port 8080 to be suppressed")
	}
}

func TestIsSuppressed_ExpiredEntry(t *testing.T) {
	base := time.Now()
	l := New()
	l.now = fixedNow(base)
	l.Suppress(8080, 1*time.Millisecond, "test")
	l.now = fixedNow(base.Add(1 * time.Second))
	if l.IsSuppressed(8080) {
		t.Fatal("expected port 8080 suppression to have expired")
	}
}

func TestIsSuppressed_UnknownPort(t *testing.T) {
	l := New()
	if l.IsSuppressed(9999) {
		t.Fatal("unknown port should not be suppressed")
	}
}

func TestRemove(t *testing.T) {
	l := New()
	l.Suppress(443, 10*time.Minute, "test")
	l.Remove(443)
	if l.IsSuppressed(443) {
		t.Fatal("expected port 443 to be unsuppressed after Remove")
	}
}

func TestExpire_RemovesStaleEntries(t *testing.T) {
	base := time.Now()
	l := New()
	l.now = fixedNow(base)
	l.Suppress(80, 1*time.Millisecond, "old")
	l.Suppress(443, 1*time.Hour, "active")
	l.now = fixedNow(base.Add(1 * time.Second))
	l.Expire()
	if l.IsSuppressed(80) {
		t.Fatal("port 80 should have been expired")
	}
	if !l.IsSuppressed(443) {
		t.Fatal("port 443 should still be suppressed")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	l := New()
	l.Suppress(22, 1*time.Hour, "ssh")
	l.Suppress(80, 1*time.Hour, "http")
	entries := l.All()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}
