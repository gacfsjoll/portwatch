package acknowledge_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/acknowledge"
)

func tempAckPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "acks", "ack.json")
}

func TestStore_AcknowledgeAndLoad(t *testing.T) {
	path := tempAckPath(t)
	s := acknowledge.NewStore(path)

	if err := s.Acknowledge(8080); err != nil {
		t.Fatalf("acknowledge: %v", err)
	}

	s2 := acknowledge.NewStore(path)
	if err := s2.Load(); err != nil {
		t.Fatalf("load: %v", err)
	}
	if !s2.IsAcknowledged(8080) {
		t.Error("expected port 8080 to be acknowledged after reload")
	}
}

func TestStore_LoadMissingFile(t *testing.T) {
	s := acknowledge.NewStore("/tmp/portwatch_nonexistent/acks.json")
	if err := s.Load(); err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
}

func TestStore_Revoke(t *testing.T) {
	path := tempAckPath(t)
	s := acknowledge.NewStore(path)
	_ = s.Acknowledge(9000)
	_ = s.Revoke(9000)
	if s.IsAcknowledged(9000) {
		t.Error("expected port 9000 to be revoked")
	}
}

func TestStore_List(t *testing.T) {
	path := tempAckPath(t)
	s := acknowledge.NewStore(path)
	_ = s.Acknowledge(80)
	_ = s.Acknowledge(443)
	list := s.List()
	if len(list) != 2 {
		t.Errorf("expected 2 ports, got %d", len(list))
	}
}

func TestStore_SaveCreatesParentDirs(t *testing.T) {
	path := tempAckPath(t)
	s := acknowledge.NewStore(path)
	if err := s.Acknowledge(22); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
