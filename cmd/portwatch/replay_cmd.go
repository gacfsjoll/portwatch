package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/example/portwatch/internal/alert"
	"github.com/example/portwatch/internal/history"
	"github.com/example/portwatch/internal/replay"
)

func runReplay(args []string) error {
	fs := flag.NewFlagSet("replay", flag.ContinueOnError)
	sinceStr := fs.String("since", "", "only replay events newer than this duration (e.g. 2h, 30m)")
	dryRun := fs.Bool("dry-run", false, "print events without dispatching notifications")
	delay := fs.Duration("delay", 0, "pause between dispatched events")

	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("replay: load config: %w", err)
	}

	var since time.Duration
	if *sinceStr != "" {
		since, err = time.ParseDuration(*sinceStr)
		if err != nil {
			return fmt.Errorf("replay: invalid --since value %q: %w", *sinceStr, err)
		}
	}

	notifier, err := alert.FromConfig(cfg.Alert)
	if err != nil {
		return fmt.Errorf("replay: build notifier: %w", err)
	}

	rec, err := history.NewRecorder(cfg.HistoryPath)
	if err != nil {
		return fmt.Errorf("replay: open history: %w", err)
	}

	opts := replay.Options{
		Since:  since,
		DryRun: *dryRun,
		Delay:  *delay,
	}

	r := replay.New(notifier, opts)
	count, err := r.Run(context.Background(), rec)
	if err != nil {
		return fmt.Errorf("replay: %w", err)
	}

	fmt.Fprintf(os.Stdout, "replayed %d event(s)\n", count)
	return nil
}
