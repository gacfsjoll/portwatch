package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/circuitbreaker"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllow_ClosedByDefault(t *testing.T) {
	b := circuitbreaker.NewWithClock(3, 30*time.Second, fixedClock(epoch))
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestAllow_OpensAfterThreshold(t *testing.T) {
	b := circuitbreaker.NewWithClock(3, 30*time.Second, fixedClock(epoch))
	b.RecordFailure()
	b.RecordFailure()
	if b.CurrentState() != circuitbreaker.StateClosed {
		t.Fatal("expected closed before threshold")
	}
	b.RecordFailure()
	if b.CurrentState() != circuitbreaker.StateOpen {
		t.Fatal("expected open after threshold")
	}
	if err := b.Allow(); err != circuitbreaker.ErrCircuitOpen {
		t.Fatalf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestAllow_RecoveryWindowReopens(t *testing.T) {
	now := epoch
	clock := func() time.Time { return now }
	b := circuitbreaker.NewWithClock(2, 10*time.Second, clock)

	b.RecordFailure()
	b.RecordFailure()
	if err := b.Allow(); err != circuitbreaker.ErrCircuitOpen {
		t.Fatal("expected open circuit")
	}

	now = epoch.Add(11 * time.Second)
	if err := b.Allow(); err != nil {
		t.Fatalf("expected circuit to allow after recovery, got %v", err)
	}
	if b.CurrentState() != circuitbreaker.StateClosed {
		t.Fatal("expected state to be closed after recovery")
	}
}

func TestRecordSuccess_ResetsClosed(t *testing.T) {
	b := circuitbreaker.NewWithClock(2, 10*time.Second, fixedClock(epoch))
	b.RecordFailure()
	b.RecordFailure()
	if b.CurrentState() != circuitbreaker.StateOpen {
		t.Fatal("expected open")
	}
	// Simulate recovery window passing then a successful probe.
	now := epoch.Add(15 * time.Second)
	b2 := circuitbreaker.NewWithClock(2, 10*time.Second, func() time.Time { return now })
	b2.RecordFailure()
	b2.RecordFailure()
	_ = b2.Allow() // triggers half-open
	b2.RecordSuccess()
	if b2.CurrentState() != circuitbreaker.StateClosed {
		t.Fatal("expected closed after success")
	}
}

func TestState_String(t *testing.T) {
	if circuitbreaker.StateClosed.String() != "closed" {
		t.Error("expected 'closed'")
	}
	if circuitbreaker.StateOpen.String() != "open" {
		t.Error("expected 'open'")
	}
}
