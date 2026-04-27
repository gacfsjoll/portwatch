// Package replay provides a mechanism to re-emit historical alert events
// through a notifier pipeline, useful for testing alert routing or recovering
// missed notifications after a daemon restart.
package replay

import (
	"context"
	"fmt"
	"time"

	"github.com/example/portwatch/internal/alert"
	"github.com/example/portwatch/internal/history"
)

// Options controls how the replay is performed.
type Options struct {
	// Since filters events older than this duration. Zero means replay all.
	Since time.Duration
	// DryRun prints events without dispatching them.
	DryRun bool
	// Delay is an optional pause between dispatched events.
	Delay time.Duration
}

// Replayer re-emits historical events through a Notifier.
type Replayer struct {
	notifier alert.Notifier
	opts     Options
}

// New returns a Replayer that dispatches events through n.
func New(n alert.Notifier, opts Options) *Replayer {
	return &Replayer{notifier: n, opts: opts}
}

// Run loads events from the recorder and dispatches them.
func (r *Replayer) Run(ctx context.Context, rec *history.Recorder) (int, error) {
	events, err := rec.Load()
	if err != nil {
		return 0, fmt.Errorf("replay: load history: %w", err)
	}

	cutoff := time.Time{}
	if r.opts.Since > 0 {
		cutoff = time.Now().UTC().Add(-r.opts.Since)
	}

	dispatched := 0
	for _, ev := range events {
		if !cutoff.IsZero() && ev.Time.Before(cutoff) {
			continue
		}
		if err := ctx.Err(); err != nil {
			return dispatched, err
		}
		if r.opts.DryRun {
			fmt.Printf("[dry-run] replay: port=%d kind=%s time=%s\n",
				ev.Port, ev.Kind, ev.Time.Format(time.RFC3339))
		} else {
			if err := r.notifier.Notify(ev); err != nil {
				return dispatched, fmt.Errorf("replay: notify port %d: %w", ev.Port, err)
			}
		}
		dispatched++
		if r.opts.Delay > 0 {
			select {
			case <-time.After(r.opts.Delay):
			case <-ctx.Done():
				return dispatched, ctx.Err()
			}
		}
	}
	return dispatched, nil
}
