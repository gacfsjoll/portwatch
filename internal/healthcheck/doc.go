// Package healthcheck provides a lightweight HTTP health endpoint for the
// portwatch daemon.
//
// It exposes a single /health route that returns a JSON Status object
// describing whether the daemon is actively scanning, how many scans have
// completed, and how long the process has been running.
//
// Usage:
//
//	srv := healthcheck.New(":9090", watchdogInstance)
//	go srv.ListenAndServe()
//
// The endpoint returns HTTP 200 when healthy and HTTP 503 when the daemon
// has not completed a scan within the configured staleness threshold.
package healthcheck
