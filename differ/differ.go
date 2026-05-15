// Package differ compares two sets of environment variables and reports
// keys that are added, removed, or changed between them.
package differ

// ChangeType describes the kind of difference detected.
type ChangeType string

const (
	Added   ChangeType = "added"
	Removed ChangeType = "removed"
	Changed ChangeType = "changed"
)

// Diff represents a single detected difference between two env maps.
type Diff struct {
	Key    string
	Type   ChangeType
	OldVal string
	NewVal string
}

// Compare returns the list of differences between a baseline env map and a
// candidate env map. Keys present only in candidate are Added; keys present
// only in baseline are Removed; keys present in both but with different
// values are Changed.
func Compare(baseline, candidate map[string]string) []Diff {
	var diffs []Diff

	// Detect removed and changed keys.
	for k, oldVal := range baseline {
		newVal, ok := candidate[k]
		if !ok {
			diffs = append(diffs, Diff{Key: k, Type: Removed, OldVal: oldVal})
			continue
		}
		if newVal != oldVal {
			diffs = append(diffs, Diff{Key: k, Type: Changed, OldVal: oldVal, NewVal: newVal})
		}
	}

	// Detect added keys.
	for k, newVal := range candidate {
		if _, ok := baseline[k]; !ok {
			diffs = append(diffs, Diff{Key: k, Type: Added, NewVal: newVal})
		}
	}

	return diffs
}

// HasChanges returns true when Compare produces at least one Diff.
func HasChanges(baseline, candidate map[string]string) bool {
	return len(Compare(baseline, candidate)) > 0
}
