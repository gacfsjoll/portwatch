package tee_test

import (
	"context"
	"errors"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/tee"
)

// stubNotifier records calls and optionally returns an error.
type stubNotifier struct {
	called int
	err    error
}

func (s *stubNotifier) Notify(_ context.Context, _ alert.Event) error {
	s.called++
	return s.err
}

func makeEvent() alert.Event {
	return alert.Event{Port: 8080, Kind: "opened"}
}

func TestNew_RequiresAtLeastOneNotifier(t *testing.T) {
	_, err := tee.New()
	if err == nil {
		t.Fatal("expected error for empty notifier list")
	}
}

func TestNotify_CallsAllNotifiers(t *testing.T) {
	a, b := &stubNotifier{}, &stubNotifier{}
	n, err := tee.New(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := n.Notify(context.Background(), makeEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.called != 1 || b.called != 1 {
		t.Errorf("expected each notifier called once, got a=%d b=%d", a.called, b.called)
	}
}

func TestNotify_ContinuesAfterFailure(t *testing.T) {
	a := &stubNotifier{err: errors.New("boom")}
	b := &stubNotifier{}
	n, _ := tee.New(a, b)
	err := n.Notify(context.Background(), makeEvent())
	if err == nil {
		t.Fatal("expected combined error")
	}
	if b.called != 1 {
		t.Errorf("second notifier should still be called, got %d", b.called)
	}
}

func TestNotify_ReturnsNilWhenAllSucceed(t *testing.T) {
	a, b := &stubNotifier{}, &stubNotifier{}
	n, _ := tee.New(a, b)
	if err := n.Notify(context.Background(), makeEvent()); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestLen_ReturnsCount(t *testing.T) {
	a, b, c := &stubNotifier{}, &stubNotifier{}, &stubNotifier{}
	n, _ := tee.New(a, b, c)
	if n.Len() != 3 {
		t.Errorf("expected 3, got %d", n.Len())
	}
}
