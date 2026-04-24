package pipeline_test

import (
	"testing"

	"github.com/user/portwatch/internal/envelope"
)

// TestEnvelopeStage_WrapsPayload verifies that a Wrapper integrated into a
// simple hand-rolled stage correctly increments sequence numbers and
// preserves the original payload so the rest of the pipeline can use it.
func TestEnvelopeStage_WrapsPayload(t *testing.T) {
	type event struct{ Port int }

	w := envelope.New()

	events := []event{{Port: 80}, {Port: 443}, {Port: 8080}}
	wrapped := make([]envelope.Envelope, 0, len(events))

	for _, ev := range events {
		wrapped = append(wrapped, w.Wrap(ev))
	}

	for i, env := range wrapped {
		expectedSeq := uint64(i + 1)
		if env.Seq != expectedSeq {
			t.Errorf("event %d: expected seq %d, got %d", i, expectedSeq, env.Seq)
		}
		got, ok := env.Payload.(event)
		if !ok {
			t.Fatalf("event %d: payload type mismatch", i)
		}
		if got.Port != events[i].Port {
			t.Errorf("event %d: expected port %d, got %d", i, events[i].Port, got.Port)
		}
	}
}

// TestEnvelopeStage_HostnameNonEmpty ensures the wrapper always populates
// the hostname field even when called from a test environment.
func TestEnvelopeStage_HostnameNonEmpty(t *testing.T) {
	w := envelope.New()
	env := w.Wrap("ping")

	if env.Hostname == "" {
		t.Fatal("hostname must not be empty")
	}
}
