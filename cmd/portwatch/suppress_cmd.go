package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/user/portwatch/internal/suppress"
)

// globalSuppressionList is shared across the daemon and CLI sub-commands
// within the same process. In a real deployment this would be persisted.
var globalSuppressionList = suppress.New()

// runSuppress handles: portwatch suppress <port> <duration> [reason]
func runSuppress(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: portwatch suppress <port> <duration> [reason]")
		os.Exit(1)
	}

	port, err := strconv.Atoi(args[0])
	if err != nil || port < 1 || port > 65535 {
		fmt.Fprintf(os.Stderr, "invalid port: %s\n", args[0])
		os.Exit(1)
	}

	d, err := time.ParseDuration(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid duration %q: %v\n", args[1], err)
		os.Exit(1)
	}

	if d <= 0 {
		fmt.Fprintf(os.Stderr, "duration must be positive, got %s\n", args[1])
		os.Exit(1)
	}

	reason := "manual suppression"
	if len(args) >= 3 {
		reason = args[2]
	}

	globalSuppressionList.Suppress(port, d, reason)
	fmt.Printf("port %d suppressed for %s (%s)\n", port, d, reason)
}

// runSuppressList prints all currently active suppression entries.
func runSuppressList() {
	entries := globalSuppressionList.All()
	if len(entries) == 0 {
		fmt.Println("no active suppressions")
		return
	}
	fmt.Printf("%-8s %-30s %s\n", "PORT", "UNTIL", "REASON")
	for _, e := range entries {
		fmt.Printf("%-8d %-30s %s\n", e.Port, e.Until.Format(time.RFC3339), e.Reason)
	}
}

// runSuppressRemove handles: portwatch suppress remove <port>
// It removes any active suppression for the given port.
func runSuppressRemove(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: portwatch suppress remove <port>")
		os.Exit(1)
	}

	port, err := strconv.Atoi(args[0])
	if err != nil || port < 1 || port > 65535 {
		fmt.Fprintf(os.Stderr, "invalid port: %s\n", args[0])
		os.Exit(1)
	}

	globalSuppressionList.Remove(port)
	fmt.Printf("suppression removed for port %d\n", port)
}
