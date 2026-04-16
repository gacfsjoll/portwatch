// Package filter provides port filtering utilities for portwatch.
// It allows including or excluding specific ports or ranges from monitoring.
package filter

import (
	"fmt"
	"strconv"
	"strings"
)

// Rule holds a parsed include or exclude filter rule.
type Rule struct {
	Low  int
	High int
}

// Filter decides which ports should be monitored.
type Filter struct {
	includes []Rule
	excludes []Rule
}

// New builds a Filter from include and exclude expression slices.
// Each expression is either a single port ("80") or a range ("8000-9000").
func New(includes, excludes []string) (*Filter, error) {
	f := &Filter{}
	for _, expr := range includes {
		r, err := parseRule(expr)
		if err != nil {
			return nil, fmt.Errorf("invalid include rule %q: %w", expr, err)
		}
		f.includes = append(f.includes, r)
	}
	for _, expr := range excludes {
		r, err := parseRule(expr)
		if err != nil {
			return nil, fmt.Errorf("invalid exclude rule %q: %w", expr, err)
		}
		f.excludes = append(f.excludes, r)
	}
	return f, nil
}

// Allow returns true if the given port should be monitored.
func (f *Filter) Allow(port int) bool {
	for _, r := range f.excludes {
		if port >= r.Low && port <= r.High {
			return false
		}
	}
	if len(f.includes) == 0 {
		return true
	}
	for _, r := range f.includes {
		if port >= r.Low && port <= r.High {
			return true
		}
	}
	return false
}

func parseRule(expr string) (Rule, error) {
	parts := strings.SplitN(strings.TrimSpace(expr), "-", 2)
	low, err := strconv.Atoi(parts[0])
	if err != nil || low < 1 || low > 65535 {
		return Rule{}, fmt.Errorf("invalid port %q", parts[0])
	}
	if len(parts) == 1 {
		return Rule{Low: low, High: low}, nil
	}
	high, err := strconv.Atoi(parts[1])
	if err != nil || high < 1 || high > 65535 || high < low {
		return Rule{}, fmt.Errorf("invalid high port %q", parts[1])
	}
	return Rule{Low: low, High: high}, nil
}
