package truncate_test

import (
	"strings"
	"testing"

	"github.com/example/portwatch/internal/truncate"
)

func TestNew_DefaultsToDefaultMaxBytes(t *testing.T) {
	tr := truncate.New(0)
	if tr.MaxBytes() != truncate.DefaultMaxBytes {
		t.Fatalf("expected %d, got %d", truncate.DefaultMaxBytes, tr.MaxBytes())
	}
}

func TestNew_NegativeDefaultsToDefaultMaxBytes(t *testing.T) {
	tr := truncate.New(-1)
	if tr.MaxBytes() != truncate.DefaultMaxBytes {
		t.Fatalf("expected %d, got %d", truncate.DefaultMaxBytes, tr.MaxBytes())
	}
}

func TestApply_ShortStringUnchanged(t *testing.T) {
	tr := truncate.New(100)
	input := "hello world"
	if got := tr.Apply(input); got != input {
		t.Fatalf("expected %q, got %q", input, got)
	}
}

func TestApply_ExactLengthUnchanged(t *testing.T) {
	tr := truncate.New(5)
	input := "hello"
	if got := tr.Apply(input); got != input {
		t.Fatalf("expected %q, got %q", input, got)
	}
}

func TestApply_LongStringTruncated(t *testing.T) {
	tr := truncate.New(50)
	input := strings.Repeat("a", 200)
	got := tr.Apply(input)
	if len(got) > 50 {
		t.Fatalf("expected len <= 50, got %d", len(got))
	}
	if !strings.HasSuffix(got, "... [truncated]") {
		t.Fatalf("expected truncation suffix, got %q", got)
	}
}

func TestApply_UTF8BoundaryRespected(t *testing.T) {
	// Each rune is 3 bytes in UTF-8; cap at 10 bytes.
	tr := truncate.New(30)
	input := strings.Repeat("日", 20) // 60 bytes total
	got := tr.Apply(input)
	if len(got) > 30 {
		t.Fatalf("output exceeds cap: len=%d", len(got))
	}
	// Result must be valid UTF-8 up to the suffix.
	withoutSuffix := strings.TrimSuffix(got, "... [truncated]")
	for i, r := range withoutSuffix {
		if r == '\uFFFD' {
			t.Fatalf("invalid UTF-8 at index %d", i)
		}
	}
}

func TestApplyf_FormatsAndTruncates(t *testing.T) {
	tr := truncate.New(20)
	got := tr.Applyf("port %d opened on host %s", 8080, "example.internal.corp.example.com")
	if len(got) > 20 {
		t.Fatalf("expected len <= 20, got %d: %q", len(got), got)
	}
}
