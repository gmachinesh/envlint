// Package patcher provides utilities for applying in-place updates to an
// existing env map — adding new keys, overwriting changed values, and
// optionally removing keys that are no longer present in the patch source.
package patcher

import "fmt"

// Op describes the kind of change applied to a single key.
type Op string

const (
	OpAdded    Op = "added"
	OpUpdated  Op = "updated"
	OpRemoved  Op = "removed"
	OpUnchanged Op = "unchanged"
)

// Change records what happened to one key during a patch operation.
type Change struct {
	Key      string
	Op       Op
	OldValue string
	NewValue string
}

// Options controls the behaviour of Patch.
type Options struct {
	// RemoveMissing deletes keys from base that are absent in the patch map.
	RemoveMissing bool
	// DryRun reports changes without mutating base.
	DryRun bool
}

// DefaultOptions returns a safe, non-destructive configuration.
func DefaultOptions() Options {
	return Options{
		RemoveMissing: false,
		DryRun:        false,
	}
}

// Patch merges patch into base according to opts and returns the list of
// changes that were (or would be, in dry-run mode) applied.
//
// base is modified in place unless DryRun is true.
func Patch(base, patch map[string]string, opts Options) []Change {
	var changes []Change

	// Apply additions and updates.
	for k, newVal := range patch {
		oldVal, exists := base[k]
		switch {
		case !exists:
			changes = append(changes, Change{Key: k, Op: OpAdded, NewValue: newVal})
			if !opts.DryRun {
				base[k] = newVal
			}
		case oldVal != newVal:
			changes = append(changes, Change{Key: k, Op: OpUpdated, OldValue: oldVal, NewValue: newVal})
			if !opts.DryRun {
				base[k] = newVal
			}
		default:
			changes = append(changes, Change{Key: k, Op: OpUnchanged, OldValue: oldVal, NewValue: newVal})
		}
	}

	// Optionally remove keys absent from the patch.
	if opts.RemoveMissing {
		for k, oldVal := range base {
			if _, ok := patch[k]; !ok {
				changes = append(changes, Change{Key: k, Op: OpRemoved, OldValue: oldVal})
				if !opts.DryRun {
					delete(base, k)
				}
			}
		}
	}

	return changes
}

// Summary returns a human-readable one-liner describing the patch outcome.
func Summary(changes []Change) string {
	var added, updated, removed, unchanged int
	for _, c := range changes {
		switch c.Op {
		case OpAdded:
			added++
		case OpUpdated:
			updated++
		case OpRemoved:
			removed++
		case OpUnchanged:
			unchanged++
		}
	}
	return fmt.Sprintf("patch: %d added, %d updated, %d removed, %d unchanged",
		added, updated, removed, unchanged)
}
