package main

import (
	"path/filepath"
	"testing"
)

func tempAckStore(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "acks.json")
}

func TestRunAcknowledge_AddsPort(t *testing.T) {
	path := tempAckStore(t)
	if err := runAcknowledge([]string{"8080"}, path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := runAcknowledgeList(path); err != nil {
		t.Fatalf("list error: %v", err)
	}
}

func TestRunAcknowledge_NoArgs(t *testing.T) {
	if err := runAcknowledge(nil, tempAckStore(t)); err == nil {
		t.Error("expected error for missing args")
	}
}

func TestRunAcknowledge_InvalidPort(t *testing.T) {
	if err := runAcknowledge([]string{"notaport"}, tempAckStore(t)); err == nil {
		t.Error("expected error for invalid port")
	}
}

func TestRunAcknowledgeRevoke(t *testing.T) {
	path := tempAckStore(t)
	_ = runAcknowledge([]string{"443"}, path)
	if err := runAcknowledgeRevoke([]string{"443"}, path); err != nil {
		t.Fatalf("revoke error: %v", err)
	}
}

func TestRunAcknowledgeRevoke_NoArgs(t *testing.T) {
	if err := runAcknowledgeRevoke(nil, tempAckStore(t)); err == nil {
		t.Error("expected error for missing args")
	}
}

func TestRunAcknowledgeRevoke_UnknownPort(t *testing.T) {
	path := tempAckStore(t)
	// Revoking a port that was never acknowledged should return an error.
	if err := runAcknowledgeRevoke([]string{"9090"}, path); err == nil {
		t.Error("expected error when revoking unacknowledged port")
	}
}

func TestParsePort_Valid(t *testing.T) {
	p, err := parsePort("80")
	if err != nil || p != 80 {
		t.Errorf("expected 80, got %d, err %v", p, err)
	}
}

func TestParsePort_Invalid(t *testing.T) {
	if _, err := parsePort("99999"); err == nil {
		t.Error("expected error for out-of-range port")
	}
}
