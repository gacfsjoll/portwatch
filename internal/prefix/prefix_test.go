package prefix_test

import (
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/prefix"
)

// stubNotifier records the last event it received and optionally returns an error.
type stubNotifier struct {
	last alert.Event
	err  error
}

func (s *stubNotifier) Notify(e alert.Event) error {
	s.last = e
	return s.err
}

func makeEvent(msg string) alert.Event {
	return alert.Event{
		Port:      8080,
		Direction: "opened",
		Message:   msg,
		Timestamp: time.Now(),
	}
}

func TestNotify_PrependsLabel(t *testing.T) {
	stub := &stubNotifier{}
	n := prefix.New("[prod]", stub)

	_ = n.Notify(makeEvent("port opened"))

	want := "[prod] port opened"
	if stub.last.Message != want {
		t.Errorf("got %q, want %q", stub.last.Message, want)
	}
}

func TestNotify_EmptyLabel_LeavesMessageUnchanged(t *testing.T) {
	stub := &stubNotifier{}
	n := prefix.New("", stub)

	_ = n.Notify(makeEvent("port closed"))

	if stub.last.Message != "port closed" {
		t.Errorf("unexpected message: %q", stub.last.Message)
	}
}

func TestNotify_PreservesOtherFields(t *testing.T) {
	stub := &stubNotifier{}
	n := prefix.New("[test]", stub)

	e := makeEvent("hello")
	e.Port = 9090
	e.Direction = "closed"
	_ = n.Notify(e)

	if stub.last.Port != 9090 {
		t.Errorf("port: got %d, want 9090", stub.last.Port)
	}
	if stub.last.Direction != "closed" {
		t.Errorf("direction: got %q, want \"closed\"", stub.last.Direction)
	}
}

func TestNotify_PropagatesUnderlyingError(t *testing.T) {
	sentinel := errors.New("backend down")
	stub := &stubNotifier{err: sentinel}
	n := prefix.New("[x]", stub)

	if err := n.Notify(makeEvent("test")); !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestNew_NilPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil notifier")
		}
	}()
	prefix.New("label", nil)
}
