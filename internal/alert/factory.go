package alert

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Config holds the configuration for building a Notifier.
type Config struct {
	// Backend selects the notifier type. Supported values: "log", "stderr", "stdout".
	Backend string
	// Output is an optional writer override (used in tests).
	Output io.Writer
}

// FromConfig constructs a Notifier from the provided Config.
// An error is returned for unknown backend values.
func FromConfig(cfg Config) (Notifier, error) {
	switch strings.ToLower(cfg.Backend) {
	case "", "log", "stderr":
		w := cfg.Output
		if w == nil {
			w = os.Stderr
		}
		return NewLogNotifier(w), nil
	case "stdout":
		w := cfg.Output
		if w == nil {
			w = os.Stdout
		}
		return NewLogNotifier(w), nil
	default:
		return nil, fmt.Errorf("alert: unknown backend %q", cfg.Backend)
	}
}

// NewPortOpenedEvent is a convenience constructor for an ALERT-level event
// signalling that a previously closed port is now open.
func NewPortOpenedEvent(port int) Event {
	return newPortEvent(LevelAlert, port, "port opened")
}

// NewPortClosedEvent is a convenience constructor for a WARN-level event
// signalling that a previously open port is now closed.
func NewPortClosedEvent(port int) Event {
	return newPortEvent(LevelWarn, port, "port closed")
}

func newPortEvent(level Level, port int, msg string) Event {
	return Event{
		Timestamp: timeNow(),
		Level:     level,
		Port:      port,
		Message:   msg,
	}
}
