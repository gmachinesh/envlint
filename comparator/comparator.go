// Package comparator compares two sets of env vars and produces a typed diff
// suitable for reporting or further processing.
package comparator

import "sort"

// ChangeKind describes the type of change detected between two env sets.
type ChangeKind string

const (
	Added    ChangeKind = "added"
	Removed  ChangeKind = "removed"
	Modified ChangeKind = "modified"
	Unchanged ChangeKind = "unchanged"
)

// Change represents a single key-level difference.
type Change struct {
	Key      string
	Kind     ChangeKind
	OldValue string
	NewValue string
}

// Result holds the full comparison output.
type Result struct {
	Changes []Change
}

// HasDiff returns true when at least one non-unchanged entry exists.
func (r Result) HasDiff() bool {
	for _, c := range r.Changes {
		if c.Kind != Unchanged {
			return true
		}
	}
	return false
}

// Compare performs a key-by-key comparison of base and head env maps.
// All keys from both maps are evaluated; order in output is alphabetical.
func Compare(base, head map[string]string) Result {
	seen := make(map[string]bool)
	var changes []Change

	for k, bv := range base {
		seen[k] = true
		if hv, ok := head[k]; !ok {
			changes = append(changes, Change{Key: k, Kind: Removed, OldValue: bv})
		} else if bv != hv {
			changes = append(changes, Change{Key: k, Kind: Modified, OldValue: bv, NewValue: hv})
		} else {
			changes = append(changes, Change{Key: k, Kind: Unchanged, OldValue: bv, NewValue: hv})
		}
	}

	for k, hv := range head {
		if !seen[k] {
			changes = append(changes, Change{Key: k, Kind: Added, NewValue: hv})
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	return Result{Changes: changes}
}
