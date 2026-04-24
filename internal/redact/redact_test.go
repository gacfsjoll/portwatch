package redact_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/redact"
)

func TestScrub_NoRules(t *testing.T) {
	r := redact.New()
	got := r.Scrub("hello world")
	if got != "hello world" {
		t.Fatalf("expected unchanged text, got %q", got)
	}
}

func TestScrub_LiteralFromConstructor(t *testing.T) {
	r := redact.New("s3cr3t")
	got := r.Scrub("token=s3cr3t; other=fine")
	want := "token=[REDACTED]; other=fine"
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestScrub_MultipleLiterals(t *testing.T) {
	r := redact.New("alpha", "beta")
	got := r.Scrub("alpha and beta are secrets")
	want := "[REDACTED] and [REDACTED] are secrets"
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestAddLiteral_EmptyIsIgnored(t *testing.T) {
	r := redact.New()
	r.AddLiteral("")
	got := r.Scrub("no change")
	if got != "no change" {
		t.Fatalf("expected unchanged text, got %q", got)
	}
}

func TestAddPattern_ValidPattern(t *testing.T) {
	r := redact.New()
	if err := r.AddPattern(`Bearer\s+\S+`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := r.Scrub("Authorization: Bearer eyJhbGciOiJIUzI1NiJ9")
	want := "Authorization: [REDACTED]"
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestAddPattern_InvalidPattern(t *testing.T) {
	r := redact.New()
	if err := r.AddPattern(`[invalid`); err == nil {
		t.Fatal("expected error for invalid regexp, got nil")
	}
}

func TestScrub_LiteralAndPatternCombined(t *testing.T) {
	r := redact.New("mysecret")
	_ = r.AddPattern(`key=[A-Za-z0-9]+`)
	got := r.Scrub("pass=mysecret key=ABCDEF12")
	want := "pass=[REDACTED] [REDACTED]"
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}
