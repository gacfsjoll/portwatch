package monitor_test

import (
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/monitor"
)

// startListener opens a TCP listener on an OS-assigned port and returns it.
func startListener(t *testing.T) (net.Listener, int) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return ln, port
}

func TestMonitor_DetectsOpenedPort(t *testing.T) {
	ln, port := startListener(t)
	ln.Close() // start closed so baseline sees it closed

	changes := make(chan monitor.Change, 4)
	m := monitor.New(port, port, 20*time.Millisecond)
	m.OnChange = func(c monitor.Change) { changes <- c }

	stop := make(chan struct{})
	errCh := make(chan error, 1)
	go func() { errCh <- m.Start(stop) }()

	// Wait for baseline scan to complete.
	time.Sleep(30 * time.Millisecond)

	// Now open the port so the next scan detects the change.
	ln2, err := net.Listen("tcp", ln.Addr().String())
	if err != nil {
		t.Skipf("could not re-open port %d: %v", port, err)
	}
	defer ln2.Close()

	select {
	case c := <-changes:
		if !c.New.Open {
			t.Errorf("expected New.Open=true, got false")
		}
		if c.Port != port {
			t.Errorf("expected port %d, got %d", port, c.Port)
		}
	case <-time.After(300 * time.Millisecond):
		t.Error("timed out waiting for change event")
	}

	close(stop)
	if err := <-errCh; err != nil {
		t.Errorf("monitor returned error: %v", err)
	}
}

func TestMonitor_StopsCleanly(t *testing.T) {
	m := monitor.New(19000, 19010, 50*time.Millisecond)
	stop := make(chan struct{})
	errCh := make(chan error, 1)
	go func() { errCh <- m.Start(stop) }()

	time.Sleep(60 * time.Millisecond)
	close(stop)

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("monitor did not stop in time")
	}
}
