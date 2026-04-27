// Package graceful provides signal-aware shutdown coordination for the
// portwatch daemon.
//
// # Overview
//
// When portwatch receives SIGINT or SIGTERM it must:
//
//  1. Stop accepting new scan cycles.
//  2. Allow any in-progress scan to complete.
//  3. Flush pending alerts and history records.
//  4. Release the PID file and health-check socket.
//
// graceful.Handler encapsulates this logic. Callers obtain a cancellable
// (or deadline-bounded) context that propagates the shutdown signal to
// every subsystem that respects context cancellation.
//
// # Usage
//
//	h := graceful.New(10 * time.Second)
//	ctx, cancel := h.Wait(context.Background())
//	defer cancel()
//	// pass ctx to monitor, health server, etc.
package graceful
