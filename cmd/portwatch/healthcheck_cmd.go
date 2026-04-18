package main

import (
	"fmt"
	"log"

	"github.com/user/portwatch/internal/healthcheck"
	"github.com/user/portwatch/internal/watchdog"
)

// runHealthServer starts the health-check HTTP server using the given watchdog
// as the health provider. It blocks until the server exits.
func runHealthServer(cfg *appConfig, wd *watchdog.Watchdog) error {
	addr := cfg.HealthAddr
	if addr == "" {
		addr = ":9090"
	}
	log.Printf("health endpoint listening on %s", addr)
	srv := healthcheck.New(addr, wd)
	return srv.ListenAndServe()
}

// printHealth performs a single health query and prints the result.
func printHealth(wd *watchdog.Watchdog) {
	if wd.Healthy() {
		fmt.Printf("status: healthy  scans: %d  last: %s\n",
			wd.ScanCount(), wd.LastScan().Format("2006-01-02 15:04:05"))
	} else {
		fmt.Println("status: unhealthy — no recent scan detected")
	}
}
