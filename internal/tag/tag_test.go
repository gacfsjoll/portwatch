package tag_test

import (
	"testing"

	"github.com/user/portwatch/internal/tag"
)

func TestNew_ValidPairs(t *testing.T) {
	s, err := tag.New([]string{"env=prod", "role=web"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 2 {
		t.Fatalf("expected 2 tags, got %d", s.Len())
	}
}

func TestNew_InvalidPair(t *testing.T) {
	_, err := tag.New([]string{"nodash"})
	if err == nil {
		t.Fatal("expected error for malformed pair")
	}
}

func TestNew_EmptyKey(t *testing.T) {
	_, err := tag.New([]string{"=value"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestGet_Present(t *testing.T) {
	s, _ := tag.New([]string{"env=staging"})
	v, ok := s.Get("env")
	if !ok || v != "staging" {
		t.Fatalf("expected env=staging, got %q ok=%v", v, ok)
	}
}

func TestGet_Missing(t *testing.T) {
	s, _ := tag.New([]string{"env=prod"})
	_, ok := s.Get("missing")
	if ok {
		t.Fatal("expected missing key to return false")
	}
}

func TestAll_ReturnsSorted(t *testing.T) {
	s, _ := tag.New([]string{"z=last", "a=first", "m=mid"})
	all := s.All()
	if all[0] != "a=first" || all[1] != "m=mid" || all[2] != "z=last" {
		t.Fatalf("unexpected order: %v", all)
	}
}

func TestString_CommaSeparated(t *testing.T) {
	s, _ := tag.New([]string{"a=1", "b=2"})
	got := s.String()
	if got != "a=1,b=2" {
		t.Fatalf("unexpected string: %q", got)
	}
}

func TestMerge_OtherOverrides(t *testing.T) {
	a, _ := tag.New([]string{"env=prod", "role=web"})
	b, _ := tag.New([]string{"env=staging", "dc=us-east"})
	m := a.Merge(b)

	env, _ := m.Get("env")
	if env != "staging" {
		t.Fatalf("expected env=staging after merge, got %q", env)
	}
	if m.Len() != 3 {
		t.Fatalf("expected 3 tags after merge, got %d", m.Len())
	}
}

func TestMerge_OriginalUnchanged(t *testing.T) {
	a, _ := tag.New([]string{"env=prod"})
	b, _ := tag.New([]string{"env=staging"})
	_ = a.Merge(b)

	v, _ := a.Get("env")
	if v != "prod" {
		t.Fatal("original set should not be mutated by Merge")
	}
}
