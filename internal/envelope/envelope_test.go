package envelope

import (
	"os"
	"strings"
	"testing"
	"time"
)

var fixedNow = func() time.Time {
	return time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
}

func TestWrap_SequenceIncrements(t *testing.T) {
	w := newWithClock(fixedNow)

	e1 := w.Wrap("first")
	e2 := w.Wrap("second")

	if e1.Seq != 1 {
		t.Fatalf("expected seq 1, got %d", e1.Seq)
	}
	if e2.Seq != 2 {
		t.Fatalf("expected seq 2, got %d", e2.Seq)
	}
}

func TestWrap_PayloadPreserved(t *testing.T) {
	w := newWithClock(fixedNow)
	payload := struct{ Port int }{Port: 8080}

	e := w.Wrap(payload)

	got, ok := e.Payload.(struct{ Port int })
	if !ok {
		t.Fatal("payload type lost")
	}
	if got.Port != 8080 {
		t.Fatalf("expected port 8080, got %d", got.Port)
	}
}

func TestWrap_TimestampIsUTC(t *testing.T) {
	w := newWithClock(fixedNow)
	e := w.Wrap(nil)

	if e.At.Location() != time.UTC {
		t.Fatalf("expected UTC, got %s", e.At.Location())
	}
	if !e.At.Equal(fixedNow()) {
		t.Fatalf("unexpected timestamp: %s", e.At)
	}
}

func TestWrap_HostnameAndPID(t *testing.T) {
	w := newWithClock(fixedNow)
	e := w.Wrap(nil)

	if e.Hostname == "" {
		t.Fatal("hostname must not be empty")
	}
	if e.PID != os.Getpid() {
		t.Fatalf("expected pid %d, got %d", os.Getpid(), e.PID)
	}
}

func TestEnvelope_String(t *testing.T) {
	w := newWithClock(fixedNow)
	e := w.Wrap("x")
	s := e.String()

	if !strings.Contains(s, "#1") {
		t.Errorf("expected seq in string, got: %s", s)
	}
	if !strings.Contains(s, "2024-06-01T12:00:00Z") {
		t.Errorf("expected timestamp in string, got: %s", s)
	}
}

func TestWrap_ConcurrentSafe(t *testing.T) {
	w := newWithClock(fixedNow)
	const n = 100
	seen := make(map[uint64]bool, n)
	results := make(chan uint64, n)

	for i := 0; i < n; i++ {
		go func() { results <- w.Wrap(nil).Seq }()
	}
	for i := 0; i < n; i++ {
		seq := <-results
		if seen[seq] {
			t.Fatalf("duplicate sequence number %d", seq)
		}
		seen[seq] = true
	}
}
