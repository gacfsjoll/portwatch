// Package graceful provides utilities for handling OS signals and
// performing an orderly shutdown of long-running daemon processes.
package graceful

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ShutdownTimeout is the default maximum time allowed for cleanup.
const ShutdownTimeout = 10 * time.Second

// Handler listens for OS termination signals and cancels a context,
// giving in-flight work a bounded window to finish.
type Handler struct {
	timeout time.Duration
	signals []os.Signal
}

// New returns a Handler that reacts to SIGINT and SIGTERM.
func New(timeout time.Duration) *Handler {
	if timeout <= 0 {
		timeout = ShutdownTimeout
	}
	return &Handler{
		timeout: timeout,
		signals: []os.Signal{syscall.SIGINT, syscall.SIGTERM},
	}
}

// Wait blocks until a signal is received, then returns a context that
// is cancelled after the configured timeout. The caller should use this
// context to coordinate shutdown of all subsystems.
func (h *Handler) Wait(parent context.Context) (context.Context, context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, h.signals...)

	ctx, cancel := context.WithCancel(parent)

	go func() {
		defer signal.Stop(ch)
		select {
		case <-ch:
		case <-parent.Done():
		}
		cancel()
	}()

	return ctx, cancel
}

// WaitWithTimeout blocks until a signal is received, then returns a
// deadline-bounded context and its cancel function.
func (h *Handler) WaitWithTimeout(parent context.Context) (context.Context, context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, h.signals...)

	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	go func() {
		defer signal.Stop(ch)
		select {
		case <-ch:
		case <-parent.Done():
		}
	}()

	ctx, cancel = context.WithTimeout(parent, h.timeout)
	return ctx, cancel
}
