package rollup_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/rollup"
)

func makeEvent(port int, kind string) alert.Event {
	return alert.Event{Port: port, Kind: kind}
}

func TestRollup_BatchesEvents(t *testing.T) {
	var mu sync.Mutex
	var received [][]alert.Event

	r := rollup.New(40*time.Millisecond, func(events []alert.Event) {
		mu.Lock()
		received = append(received, events)
		mu.Unlock()
	})

	r.Add(makeEvent(80, "opened"))
	r.Add(makeEvent(443, "opened"))
	r.Add(makeEvent(22, "closed"))

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if len(received) != 1 {
		t.Fatalf("expected 1 flush, got %d", len(received))
	}
	if len(received[0]) != 3 {
		t.Errorf("expected 3 events in batch, got %d", len(received[0]))
	}
}

func TestRollup_TimerResets(t *testing.T) {
	var mu sync.Mutex
	flushCount := 0

	r := rollup.New(60*time.Millisecond, func(_ []alert.Event) {
		mu.Lock()
		flushCount++
		mu.Unlock()
	})

	// Add events spaced closer than the window; only one flush should occur.
	for i := 0; i < 4; i++ {
		r.Add(makeEvent(8000+i, "opened"))
		time.Sleep(20 * time.Millisecond)
	}
	time.Sleep(120 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if flushCount != 1 {
		t.Errorf("expected 1 flush, got %d", flushCount)
	}
}

func TestRollup_ForceFlush(t *testing.T) {
	var mu sync.Mutex
	var received [][]alert.Event

	r := rollup.New(10*time.Second, func(events []alert.Event) {
		mu.Lock()
		received = append(received, events)
		mu.Unlock()
	})

	r.Add(makeEvent(9090, "opened"))
	r.Flush() // should not wait for the 10 s window

	mu.Lock()
	defer mu.Unlock()
	if len(received) != 1 || len(received[0]) != 1 {
		t.Errorf("expected 1 batch with 1 event after forced flush")
	}
}

func TestRollup_EmptyFlushIsNoop(t *testing.T) {
	called := false
	r := rollup.New(20*time.Millisecond, func(_ []alert.Event) {
		called = true
	})
	r.Flush()
	if called {
		t.Error("flush on empty roller should not invoke callback")
	}
}
