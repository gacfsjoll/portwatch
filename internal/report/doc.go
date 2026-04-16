// Package report provides functionality for generating human-readable
// summaries of port monitoring activity recorded by the history package.
//
// Usage:
//
//	events, _ := history.Load("/var/lib/portwatch/history.json")
//	summary := report.FromHistory(events)
//	g := report.NewGenerator(os.Stdout)
//	g.Print(summary)
package report
