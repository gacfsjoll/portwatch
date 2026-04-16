// Package report generates summaries of port monitoring activity.
package report

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/history"
)

// Summary holds aggregated report data.
type Summary struct {
	GeneratedAt time.Time
	TotalEvents int
	Opened      []string
	Closed      []string
}

// Generator builds reports from recorded history.
type Generator struct {
	Out io.Writer
}

// NewGenerator returns a Generator writing to w. If w is nil, os.Stdout is used.
func NewGenerator(w io.Writer) *Generator {
	if w == nil {
		w = os.Stdout
	}
	return &Generator{Out: w}
}

// FromHistory builds a Summary from the given events.
func FromHistory(events []history.Event) Summary {
	s := Summary{GeneratedAt: time.Now()}
	for _, e := range events {
		s.TotalEvents++
		switch e.Kind {
		case "opened":
			s.Opened = append(s.Opened, fmt.Sprintf("%d/%s", e.Port, e.Proto))
		case "closed":
			s.Closed = append(s.Closed, fmt.Sprintf("%d/%s", e.Port, e.Proto))
		}
	}
	return s
}

// Print writes a human-readable summary to g.Out.
func (g *Generator) Print(s Summary) {
	fmt.Fprintf(g.Out, "Port Watch Report — %s\n", s.GeneratedAt.Format(time.RFC1123))
	fmt.Fprintf(g.Out, "Total events : %d\n", s.TotalEvents)
	fmt.Fprintf(g.Out, "Ports opened : %d\n", len(s.Opened))
	for _, p := range s.Opened {
		fmt.Fprintf(g.Out, "  + %s\n", p)
	}
	fmt.Fprintf(g.Out, "Ports closed : %d\n", len(s.Closed))
	for _, p := range s.Closed {
		fmt.Fprintf(g.Out, "  - %s\n", p)
	}
}
