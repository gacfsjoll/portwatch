package deadman_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/deadman"
)

func fixedClock(t time.Time) func() time.Time {
	var mu time.Time = t
	return func() time.Time { return mu }
}

func TestSwitch_DoesNotFireWhenResetCalled(t *testing.T) {
	var fired int32
	notify := func(_ time.Time, _ time.Duration) { atomic.AddInt32(&fired, 1) }

	sw := deadman.New(200*time.Millisecond, notify)
	defer sw.Stop()

	for i := 0; i < 5; i++ {
		time.Sleep(50 * time.Millisecond)
		sw.Reset()
	}

	if atomic.LoadInt32(&fired) != 0 {
		t.Fatalf("expected no alerts, got %d", fired)
	}
}

func TestSwitch_FiresWhenDeadlineExceeded(t *testing.T) {
	var fired int32
	notify := func(_ time.Time, _ time.Duration) { atomic.AddInt32(&fired, 1) }

	sw := deadman.New(100*time.Millisecond, notify)
	defer sw.Stop()

	// Do NOT call Reset — let the deadline expire.
	time.Sleep(300 * time.Millisecond)

	if atomic.LoadInt32(&fired) == 0 {
		t.Fatal("expected alert to fire, but it did not")
	}
}

func TestSwitch_StopsCleanly(t *testing.T) {
	notify := func(_ time.Time, _ time.Duration) {}
	sw := deadman.New(500*time.Millisecond, notify)

	done := make(chan struct{})
	go func() {
		sw.Stop()
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(2 * time.Second):
		t.Fatal("Stop() did not return in time")
	}
}

func TestSwitch_LastSeen_UpdatesOnReset(t *testing.T) {
	notify := func(_ time.Time, _ time.Duration) {}
	sw := deadman.New(5*time.Second, notify)
	defer sw.Stop()

	before := sw.LastSeen()
	time.Sleep(10 * time.Millisecond)
	sw.Reset()
	after := sw.LastSeen()

	if !after.After(before) {
		t.Errorf("expected LastSeen to advance after Reset; before=%v after=%v", before, after)
	}
}
