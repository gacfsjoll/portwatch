// Package healthcheck exposes a simple HTTP endpoint that reports daemon health.
package healthcheck

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Status holds the current health state of the daemon.
type Status struct {
	Healthy    bool      `json:"healthy"`
	LastScan   time.Time `json:"last_scan,omitempty"`
	ScanCount  int64     `json:"scan_count"`
	Uptime     string    `json:"uptime"`
	Message    string    `json:"message,omitempty"`
}

// Provider is the interface required to retrieve health information.
type Provider interface {
	Healthy() bool
	LastScan() time.Time
	ScanCount() int64
}

// Server serves the health endpoint.
type Server struct {
	provider  Provider
	started   time.Time
	addr      string
}

// New creates a new health check Server.
func New(addr string, p Provider) *Server {
	return &Server{provider: p, addr: addr, started: time.Now()}
}

// ListenAndServe starts the HTTP server. It blocks until the server stops.
func (s *Server) ListenAndServe() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	return http.ListenAndServe(s.addr, mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	st := Status{
		Healthy:   s.provider.Healthy(),
		LastScan:  s.provider.LastScan(),
		ScanCount: s.provider.ScanCount(),
		Uptime:    fmt.Sprintf("%.0fs", time.Since(s.started).Seconds()),
	}
	if !st.Healthy {
		st.Message = "no recent scan detected"
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(st)
}
