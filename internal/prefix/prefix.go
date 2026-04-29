// Package prefix provides a simple string-prefix notifier decorator that
// prepends a fixed label to every alert message before forwarding it to
// an underlying Notifier.  This is useful when multiple portwatch instances
// share the same alerting backend and need to be distinguished by host or
// environment name.
package prefix

import (
	"fmt"

	"github.com/user/portwatch/internal/alert"
)

// Notifier wraps another alert.Notifier and prepends Label to every message.
type Notifier struct {
	Label string
	next  alert.Notifier
}

// New returns a Notifier that prepends label (e.g. "[prod]") to every event
// message before delegating to next.  If next is nil, New panics.
func New(label string, next alert.Notifier) *Notifier {
	if next == nil {
		panic("prefix: underlying notifier must not be nil")
	}
	return &Notifier{Label: label, next: next}
}

// Notify prepends the configured label to e.Message and forwards the modified
// event to the underlying notifier.
func (n *Notifier) Notify(e alert.Event) error {
	copy := e
	if n.Label != "" {
		copy.Message = fmt.Sprintf("%s %s", n.Label, e.Message)
	}
	return n.next.Notify(copy)
}
