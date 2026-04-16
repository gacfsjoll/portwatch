package main

import (
	"fmt"
	"os"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/report"
)

// runReport loads recorded history and prints a summary to stdout.
func runReport(historyPath string) error {
	events, err := history.Load(historyPath)
	if err != nil {
		return fmt.Errorf("loading history: %w", err)
	}
	if len(events) == 0 {
		fmt.Fprintln(os.Stdout, "No events recorded yet.")
		return nil
	}
	g := report.NewGenerator(os.Stdout)
	s := report.FromHistory(events)
	g.Print(s)
	return nil
}
