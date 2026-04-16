package snapshot_test

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func ports(nums ...int) []scanner.PortState {
	states := make([]scanner.PortState, len(nums))
	for i, n := range nums {
		states[i] = scanner.PortState{Port: n, Open: true}
	}
	return states
}

func TestCompare_NoChanges(t *testing.T) {
	diff := snapshot.Compare(ports(80, 443), ports(80, 443))
	if diff.HasChanges() {
		t.Errorf("expected no changes, got opened=%v closed=%v", diff.Opened, diff.Closed)
	}
}

func TestCompare_DetectsOpened(t *testing.T) {
	diff := snapshot.Compare(ports(80), ports(80, 8080))
	if len(diff.Opened) != 1 || diff.Opened[0].Port != 8080 {
		t.Errorf("expected port 8080 opened, got %v", diff.Opened)
	}
	if len(diff.Closed) != 0 {
		t.Errorf("expected no closed ports, got %v", diff.Closed)
	}
}

func TestCompare_DetectsClosed(t *testing.T) {
	diff := snapshot.Compare(ports(80, 443), ports(80))
	if len(diff.Closed) != 1 || diff.Closed[0].Port != 443 {
		t.Errorf("expected port 443 closed, got %v", diff.Closed)
	}
	if len(diff.Opened) != 0 {
		t.Errorf("expected no opened ports, got %v", diff.Opened)
	}
}

func TestCompare_EmptyPrevious(t *testing.T) {
	diff := snapshot.Compare(nil, ports(22, 80))
	if len(diff.Opened) != 2 {
		t.Errorf("expected 2 opened ports, got %v", diff.Opened)
	}
}

func TestCompare_EmptyCurrent(t *testing.T) {
	diff := snapshot.Compare(ports(22, 80), nil)
	if len(diff.Closed) != 2 {
		t.Errorf("expected 2 closed ports, got %v", diff.Closed)
	}
}
