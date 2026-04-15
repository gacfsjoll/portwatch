package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
)

func TestFromConfig_DefaultsToLogNotifier(t *testing.T) {
	var buf bytes.Buffer
	n, err := alert.FromConfig(alert.Config{Backend: "", Output: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestFromConfig_StdoutBackend(t *testing.T) {
	var buf bytes.Buffer
	n, err := alert.FromConfig(alert.Config{Backend: "stdout", Output: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e := makeEvent(alert.LevelInfo, 443, "https")
	if err := n.Notify(e); err != nil {
		t.Fatalf("notify error: %v", err)
	}
	if !strings.Contains(buf.String(), "port=443") {
		t.Errorf("expected port=443 in output, got: %s", buf.String())
	}
}

func TestFromConfig_UnknownBackend(t *testing.T) {
	_, err := alert.FromConfig(alert.Config{Backend: "slack"})
	if err == nil {
		t.Fatal("expected error for unknown backend")
	}
	if !strings.Contains(err.Error(), "slack") {
		t.Errorf("error message should mention backend name, got: %v", err)
	}
}

func TestNewPortOpenedEvent(t *testing.T) {
	e := alert.NewPortOpenedEvent(3000)
	if e.Port != 3000 {
		t.Errorf("expected port 3000, got %d", e.Port)
	}
	if e.Level != alert.LevelAlert {
		t.Errorf("expected ALERT level, got %s", e.Level)
	}
}

func TestNewPortClosedEvent(t *testing.T) {
	e := alert.NewPortClosedEvent(3000)
	if e.Port != 3000 {
		t.Errorf("expected port 3000, got %d", e.Port)
	}
	if e.Level != alert.LevelWarn {
		t.Errorf("expected WARN level, got %s", e.Level)
	}
}
