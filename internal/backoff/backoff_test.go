package backoff

import (
	"testing"
	"time"
)

func TestNext_FirstAttemptReturnsInitial(t *testing.T) {
	p := Policy{Initial: 1 * time.Second, Max: 60 * time.Second, Multiplier: 2.0}
	b := New(p)

	d := b.Next("scanner")
	if d != 1*time.Second {
		t.Fatalf("expected 1s, got %v", d)
	}
}

func TestNext_DoublesEachAttempt(t *testing.T) {
	p := Policy{Initial: 1 * time.Second, Max: 60 * time.Second, Multiplier: 2.0}
	b := New(p)

	b.Next("k") // attempt 0 -> 1s
	d := b.Next("k") // attempt 1 -> 2s
	if d != 2*time.Second {
		t.Fatalf("expected 2s, got %v", d)
	}
}

func TestNext_CapsAtMax(t *testing.T) {
	p := Policy{Initial: 1 * time.Second, Max: 5 * time.Second, Multiplier: 2.0}
	b := New(p)

	var last time.Duration
	for i := 0; i < 10; i++ {
		last = b.Next("k")
	}
	if last > 5*time.Second {
		t.Fatalf("expected delay capped at 5s, got %v", last)
	}
}

func TestReset_ClearsAttempts(t *testing.T) {
	p := DefaultPolicy()
	b := New(p)

	b.Next("host")
	b.Next("host")
	if b.Attempts("host") != 2 {
		t.Fatalf("expected 2 attempts before reset")
	}

	b.Reset("host")
	if b.Attempts("host") != 0 {
		t.Fatalf("expected 0 attempts after reset, got %d", b.Attempts("host"))
	}
}

func TestNext_IndependentKeys(t *testing.T) {
	p := Policy{Initial: 1 * time.Second, Max: 60 * time.Second, Multiplier: 2.0}
	b := New(p)

	b.Next("a")
	b.Next("a")
	d := b.Next("b") // first attempt for "b"
	if d != 1*time.Second {
		t.Fatalf("expected 1s for independent key, got %v", d)
	}
}

func TestDefaultPolicy_MultiplierFallback(t *testing.T) {
	p := Policy{Initial: 200 * time.Millisecond, Max: 10 * time.Second, Multiplier: 0}
	b := NewWithClock(p, time.Now)
	// multiplier should default to 2.0
	b.Next("x")
	d := b.Next("x")
	if d < 200*time.Millisecond {
		t.Fatalf("expected delay >= initial, got %v", d)
	}
}

func TestAttempts_TracksCorrectly(t *testing.T) {
	b := New(DefaultPolicy())
	for i := 0; i < 5; i++ {
		b.Next("port:8080")
	}
	if got := b.Attempts("port:8080"); got != 5 {
		t.Fatalf("expected 5 attempts, got %d", got)
	}
}
