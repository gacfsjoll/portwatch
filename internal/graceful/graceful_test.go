package graceful_test

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/user/portwatch/internal/graceful"
)

func TestNew_DefaultTimeout(t *testing.T) {
	h := graceful.New(0)
	if h == nil {
		t.Fatal("expected non-nil Handler")
	}
}

func TestNew_CustomTimeout(t *testing.T) {
	h := graceful.New(5 * time.Second)
	if h == nil {
		t.Fatal("expected non-nil Handler")
	}
}

func TestWait_CancelledParentExits(t *testing.T) {
	h := graceful.New(graceful.ShutdownTimeout)
	parent, parentCancel := context.WithCancel(context.Background())

	ctx, cancel := h.Wait(parent)
	defer cancel()

	// Cancelling the parent should propagate to the returned context.
	parentCancel()

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(time.Second):
		t.Fatal("context was not cancelled after parent cancel")
	}
}

func TestWait_SignalCancelsContext(t *testing.T) {
	h := graceful.New(graceful.ShutdownTimeout)
	ctx, cancel := h.Wait(context.Background())
	defer cancel()

	// Send SIGINT to ourselves.
	syscall.Kill(syscall.Getpid(), syscall.SIGINT) //nolint:errcheck

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(2 * time.Second):
		t.Fatal("context was not cancelled after SIGINT")
	}
}

func TestWaitWithTimeout_ReturnsDeadlineContext(t *testing.T) {
	h := graceful.New(50 * time.Millisecond)
	parent, parentCancel := context.WithCancel(context.Background())
	defer parentCancel()

	ctx, cancel := h.WaitWithTimeout(parent)
	defer cancel()

	if _, ok := ctx.Deadline(); !ok {
		t.Fatal("expected context to have a deadline")
	}
}
