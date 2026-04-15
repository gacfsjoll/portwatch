package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestDefault(t *testing.T) {
	cfg := config.Default()
	if cfg.Scan.PortStart != 1 {
		t.Errorf("expected PortStart=1, got %d", cfg.Scan.PortStart)
	}
	if cfg.Scan.PortEnd != 65535 {
		t.Errorf("expected PortEnd=65535, got %d", cfg.Scan.PortEnd)
	}
	if cfg.Scan.Interval != 30*time.Second {
		t.Errorf("expected Interval=30s, got %v", cfg.Scan.Interval)
	}
	if cfg.Alert.Backend != "log" {
		t.Errorf("expected Backend=log, got %q", cfg.Alert.Backend)
	}
}

func TestLoad_ValidFile(t *testing.T) {
	path := writeTempConfig(t, `
scan:
  port_start: 1024
  port_end: 9000
  interval: 10s
alert:
  backend: stdout
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Scan.PortStart != 1024 {
		t.Errorf("expected 1024, got %d", cfg.Scan.PortStart)
	}
	if cfg.Alert.Backend != "stdout" {
		t.Errorf("expected stdout, got %q", cfg.Alert.Backend)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestValidate_InvalidRange(t *testing.T) {
	cfg := config.Default()
	cfg.Scan.PortStart = 9000
	cfg.Scan.PortEnd = 1000
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for inverted port range")
	}
}

func TestValidate_ZeroInterval(t *testing.T) {
	cfg := config.Default()
	cfg.Scan.Interval = 0
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for zero interval")
	}
}
