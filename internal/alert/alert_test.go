package alert_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func makeEvent(level alert.Level, port int, msg string) alert.Event {
	return alert.Event{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Level:     level,
		Port:      port,
		Message:   msg,
	}
}

func TestLogNotifier_Notify(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewLogNotifier(&buf)

	e := makeEvent(alert.LevelAlert, 8080, "unexpected port opened")
	if err := n.Notify(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"ALERT", "port=8080", "unexpected port opened", "2024-01-15"} {
		if !strings.Contains(out, want) {
			t.Errorf("output %q missing %q", out, want)
		}
	}
}

func TestLogNotifier_DefaultsToStderr(t *testing.T) {
	n := alert.NewLogNotifier(nil)
	if n.Out == nil {
		t.Fatal("expected non-nil Out when nil passed")
	}
}

func TestMultiNotifier_Notify(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	n1 := alert.NewLogNotifier(&buf1)
	n2 := alert.NewLogNotifier(&buf2)
	multi := alert.NewMultiNotifier(n1, n2)

	e := makeEvent(alert.LevelWarn, 22, "ssh port detected")
	if err := multi.Notify(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf1.String(), "port=22") {
		t.Errorf("notifier 1 did not receive event")
	}
	if !strings.Contains(buf2.String(), "port=22") {
		t.Errorf("notifier 2 did not receive event")
	}
}

func TestMultiNotifier_Empty(t *testing.T) {
	multi := alert.NewMultiNotifier()
	e := makeEvent(alert.LevelInfo, 80, "http")
	if err := multi.Notify(e); err != nil {
		t.Fatalf("unexpected error on empty MultiNotifier: %v", err)
	}
}
