package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/user/portwatch/internal/audit"
)

func runAuditLog(cfg interface{ AuditLogPath() string }) error {
	path := cfg.AuditLogPath()
	entries, err := audit.Load(path)
	if err != nil {
		return fmt.Errorf("audit: %w", err)
	}
	if len(entries) == 0 {
		fmt.Println("no audit entries found")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tACTOR\tACTION\tPORT\tDETAIL")
	for _, e := range entries {
		port := "-"
		if e.Port != 0 {
			port = fmt.Sprintf("%d", e.Port)
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			e.Timestamp.Format("2006-01-02T15:04:05Z"),
			e.Actor,
			e.Action,
			port,
			e.Detail,
		)
	}
	return w.Flush()
}
