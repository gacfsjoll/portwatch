// Package monitor provides a polling daemon that wraps the scanner package
// to detect and report changes in port state over time.
//
// Basic usage:
//
//	m := monitor.New(1024, 65535, 5*time.Second)
//	m.OnChange = func(c monitor.Change) {
//		fmt.Printf("port %d changed: %s -> %s\n", c.Port, c.Old, c.New)
//	}
//
//	stop := make(chan struct{})
//	signal.Notify(...) // wire up OS signals to close stop
//	if err := m.Start(stop); err != nil {
//		log.Fatal(err)
//	}
//
// The OnChange callback is invoked synchronously inside the scan loop;
// keep it non-blocking or dispatch to a goroutine if heavy work is needed.
package monitor
