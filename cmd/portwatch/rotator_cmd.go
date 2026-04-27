package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"
)

// runRotatorInfo prints information about existing rotated log backups found
// next to the configured log path.
func runRotatorInfo(args []string) error {
	fs := flag.NewFlagSet("rotator-info", flag.ContinueOnError)
	logPath := fs.String("log", "", "path to the portwatch log file")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *logPath == "" {
		return fmt.Errorf("--log is required")
	}

	pattern := *logPath + ".*"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("glob %q: %w", pattern, err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "BACKUP\tSIZE (bytes)\tMODIFIED")
	for _, m := range matches {
		fi, err := os.Stat(m)
		if err != nil {
			continue
		}
		fmt.Fprintf(w, "%s\t%d\t%s\n",
			filepath.Base(m),
			fi.Size(),
			fi.ModTime().UTC().Format(time.RFC3339),
		)
	}
	_ = w.Flush()

	if len(matches) == 0 {
		fmt.Fprintln(os.Stdout, "(no rotated backups found)")
	}
	return nil
}

// runRotatorPrune removes rotated log backups older than the given duration.
func runRotatorPrune(args []string) error {
	fs := flag.NewFlagSet("rotator-prune", flag.ContinueOnError)
	logPath := fs.String("log", "", "path to the portwatch log file")
	older := fs.Duration("older-than", 7*24*time.Hour, "remove backups older than this duration")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *logPath == "" {
		return fmt.Errorf("--log is required")
	}

	cutoff := time.Now().Add(-*older)
	pattern := *logPath + ".*"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("glob %q: %w", pattern, err)
	}

	removed := 0
	for _, m := range matches {
		fi, err := os.Stat(m)
		if err != nil {
			continue
		}
		if fi.ModTime().Before(cutoff) {
			if err := os.Remove(m); err == nil {
				fmt.Fprintf(os.Stdout, "removed %s\n", filepath.Base(m))
				removed++
			}
		}
	}
	fmt.Fprintf(os.Stdout, "%d backup(s) pruned\n", removed)
	return nil
}
