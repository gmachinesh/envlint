// Package filter provides utilities for filtering validation results
// based on severity levels and variable name patterns.
package filter

import (
	"strings"

	"github.com/yourorg/envlint/validator"
)

// Options holds the filtering criteria for validation results.
type Options struct {
	// Severity filters results to only include those at or above this level.
	// Accepted values: "error", "warning", "info". Empty string means no filter.
	Severity string

	// Prefix filters results to only include variables with this prefix.
	// Empty string means no filter.
	Prefix string
}

// Apply filters a slice of validation results according to the given Options.
// It returns a new slice containing only the results that match all criteria.
func Apply(results []validator.Result, opts Options) []validator.Result {
	filtered := make([]validator.Result, 0, len(results))
	for _, r := range results {
		if !matchesSeverity(r, opts.Severity) {
			continue
		}
		if !matchesPrefix(r, opts.Prefix) {
			continue
		}
		filtered = append(filtered, r)
	}
	return filtered
}

// matchesSeverity returns true if the result's severity is at or above the
// requested level. If the requested level is empty, all results match.
func matchesSeverity(r validator.Result, severity string) bool {
	if severity == "" {
		return true
	}
	order := map[string]int{
		"info":    0,
		"warning": 1,
		"error":   2,
	}
	want, ok := order[strings.ToLower(severity)]
	if !ok {
		return true
	}
	got, ok := order[strings.ToLower(r.Severity)]
	if !ok {
		return true
	}
	return got >= want
}

// matchesPrefix returns true if the result's variable name starts with the
// given prefix. If prefix is empty, all results match.
func matchesPrefix(r validator.Result, prefix string) bool {
	if prefix == "" {
		return true
	}
	return strings.HasPrefix(r.Variable, prefix)
}
