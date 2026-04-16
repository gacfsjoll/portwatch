// Package snapshot provides utilities for comparing two sets of scanned port
// states and determining which ports have been opened or closed between scans.
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
//	}
package snapshot
