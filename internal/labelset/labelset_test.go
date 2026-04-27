package labelset_test

import (
	"testing"

	"github.com/user/portwatch/internal/labelset"
)

func TestNew_ValidPairs(t *testing.T) {
	ls, err := labelset.New("env=prod", "owner=alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := ls.Get("env"); !ok || v != "prod" {
		t.Errorf("expected env=prod, got %q ok=%v", v, ok)
	}
}

func TestNew_InvalidPair(t *testing.T) {
	_, err := labelset.New("noequalssign")
	if err == nil {
		t.Fatal("expected error for pair without '='")
	}
}

func TestNew_EmptyKey(t *testing.T) {
	_, err := labelset.New("=value")
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestGet_MissingKey(t *testing.T) {
	ls, _ := labelset.New("a=1")
	_, ok := ls.Get("missing")
	if ok {
		t.Error("expected ok=false for missing key")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	ls, _ := labelset.New("x=1", "y=2")
	m := ls.All()
	m["x"] = "mutated"
	if v, _ := ls.Get("x"); v != "1" {
		t.Error("All() should return a copy, not expose internal map")
	}
}

func TestString_IsDeterministic(t *testing.T) {
	ls, _ := labelset.New("z=3", "a=1", "m=2")
	got := ls.String()
	want := "a=1,m=2,z=3"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestMerge_OtherOverwrites(t *testing.T) {
	a, _ := labelset.New("env=staging", "owner=alice")
	b, _ := labelset.New("env=prod")
	merged := a.Merge(b)
	if v, _ := merged.Get("env"); v != "prod" {
		t.Errorf("expected merged env=prod, got %q", v)
	}
	if v, _ := merged.Get("owner"); v != "alice" {
		t.Errorf("expected owner=alice to survive merge, got %q", v)
	}
}

func TestMerge_OriginalUnchanged(t *testing.T) {
	a, _ := labelset.New("env=staging")
	b, _ := labelset.New("env=prod")
	a.Merge(b)
	if v, _ := a.Get("env"); v != "staging" {
		t.Error("Merge must not modify the receiver")
	}
}
