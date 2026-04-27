// Package snapshot provides utilities for comparing two sets of scanned port
// states and determining which ports have been opened or closed between scans.
//
// The primary entry point is [Compare], which accepts two slices of
// [scanner.PortState] values representing a previous and current scan, and
// returns a [Diff] describing what changed between them.
//
// Usage:
//
//	prev := []scanner.PortState{{Port: 80, Open: true}}
//	curr := []scanner.PortState{{Port: 80, Open: true}, {Port: 8080, Open: true}}
//
//	diff := snapshot.Compare(prev, curr)
//	if diff.HasChanges() {
//		for _, p := range diff.Opened {
//			fmt.Println("opened:", p.Port)
//		}
//		for _, p := range diff.Closed {
//			fmt.Println("closed:", p.Port)
//		}
//	}
package snapshot
