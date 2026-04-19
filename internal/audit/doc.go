// Package audit provides an append-only structured audit log for portwatch.
//
// Each action taken by the daemon or CLI (such as acknowledging a port,
// adding a suppression rule, or revoking an acknowledgement) is recorded
// as a JSON entry in a newline-delimited log file.
//
// # Usage
//
//	logger := audit.New("/var/lib/portwatch/audit.jsonl")
//	logger.Log("cli", "acknowledge", 8080, "added by operator")
//
//	entries, err := audit.Load("/var/lib/portwatch/audit.jsonl")
//
// Entries are safe to write concurrently.
package audit
