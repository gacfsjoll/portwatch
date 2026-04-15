// Package alert provides notification mechanisms for port change events.
package alert

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event describes a port change that triggered an alert.
type Event struct {
	Timestamp time.Time
	Level     Level
	Port      int
	Message   string
}

// Notifier is the interface implemented by all alert backends.
type Notifier interface {
	Notify(e Event) error
}

// LogNotifier writes alerts as formatted lines to an io.Writer.
type LogNotifier struct {
	Out io.Writer
}

// NewLogNotifier returns a LogNotifier that writes to w.
// If w is nil, os.Stderr is used.
func NewLogNotifier(w io.Writer) *LogNotifier {
	if w == nil {
		w = os.Stderr
	}
	return &LogNotifier{Out: w}
}

// Notify formats the event and writes it to the configured writer.
func (l *LogNotifier) Notify(e Event) error {
	_, err := fmt.Fprintf(
		l.Out,
		"[%s] %s port=%d msg=%q\n",
		e.Timestamp.Format(time.RFC3339),
		e.Level,
		e.Port,
		e.Message,
	)
	return err
}

// MultiNotifier fans out a single event to multiple Notifiers.
type MultiNotifier struct {
	notifiers []Notifier
}

// NewMultiNotifier returns a MultiNotifier wrapping the provided notifiers.
func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{notifiers: notifiers}
}

// Notify calls every underlying Notifier and returns the first error encountered.
func (m *MultiNotifier) Notify(e Event) error {
	for _, n := range m.notifiers {
		if err := n.Notify(e); err != nil {
			return err
		}
	}
	return nil
}
