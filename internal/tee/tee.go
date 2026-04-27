// Package tee provides a fan-out notifier that duplicates alert events
// to multiple independent notification pipelines without short-circuiting
// on individual failures.
package tee

import (
	"context"
	"fmt"
	"strings"

	"github.com/user/portwatch/internal/alert"
)

// Notifier fans out a single alert.Event to a list of downstream notifiers.
// Each notifier is called in order; errors are collected and returned as a
// combined error so that a failing backend never silences the others.
type Notifier struct {
	notifiers []alert.Notifier
}

// New returns a Notifier that forwards every event to each of the supplied
// notifiers. At least one notifier must be provided.
func New(notifiers ...alert.Notifier) (*Notifier, error) {
	if len(notifiers) == 0 {
		return nil, fmt.Errorf("tee: at least one notifier is required")
	}
	return &Notifier{notifiers: notifiers}, nil
}

// Notify delivers ev to every downstream notifier. All notifiers are invoked
// regardless of intermediate failures. If one or more notifiers return an
// error the combined message is returned; otherwise nil is returned.
func (t *Notifier) Notify(ctx context.Context, ev alert.Event) error {
	var errs []string
	for _, n := range t.notifiers {
		if err := n.Notify(ctx, ev); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("tee: %d notifier(s) failed: %s", len(errs), strings.Join(errs, "; "))
	}
	return nil
}

// Len returns the number of downstream notifiers.
func (t *Notifier) Len() int { return len(t.notifiers) }
