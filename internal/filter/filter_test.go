package filter_test

import (
	"testing"

	"github.com/user/portwatch/internal/filter"
)

func TestNew_InvalidInclude(t *testing.T) {
	_, err := filter.New([]string{"abc"}, nil)
	if err == nil {
		t.Fatal("expected error for invalid include rule")
	}
}

func TestNew_InvalidExclude(t *testing.T) {
	_, err := filter.New(nil, []string{"70000"})
	if err == nil {
		t.Fatal("expected error for out-of-range exclude rule")
	}
}

func TestAllow_NoRules(t *testing.T) {
	f, err := filter.New(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, port := range []int{1, 80, 443, 65535} {
		if !f.Allow(port) {
			t.Errorf("expected port %d to be allowed with no rules", port)
		}
	}
}

func TestAllow_IncludeRange(t *testing.T) {
	f, err := filter.New([]string{"8000-9000"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !f.Allow(8080) {
		t.Error("expected 8080 to be allowed")
	}
	if f.Allow(443) {
		t.Error("expected 443 to be denied")
	}
}

func TestAllow_ExcludeSingle(t *testing.T) {
	f, err := filter.New(nil, []string{"22"})
	if err != nil {
		t.Fatal(err)
	}
	if f.Allow(22) {
		t.Error("expected port 22 to be excluded")
	}
	if !f.Allow(80) {
		t.Error("expected port 80 to be allowed")
	}
}

func TestAllow_ExcludeOverridesInclude(t *testing.T) {
	f, err := filter.New([]string{"80-90"}, []string{"85"})
	if err != nil {
		t.Fatal(err)
	}
	if !f.Allow(80) {
		t.Error("expected 80 to be allowed")
	}
	if f.Allow(85) {
		t.Error("expected 85 to be excluded despite include range")
	}
}

func TestNew_InvalidRange_HighLessThanLow(t *testing.T) {
	_, err := filter.New([]string{"9000-8000"}, nil)
	if err == nil {
		t.Fatal("expected error when high < low")
	}
}
