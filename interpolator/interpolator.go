// Package interpolator resolves cross-variable references within a .env map,
// allowing values like FOO=${BAR}_suffix to be fully expanded in place.
package interpolator

import (
	"fmt"
	"strings"

	"github.com/user/envlint/expander"
)

// Options controls interpolation behaviour.
type Options struct {
	// MaxPasses limits the number of resolution rounds to prevent infinite loops.
	MaxPasses int
	// FailOnUnresolved causes Interpolate to return an error when a reference
	// cannot be resolved after all passes.
	FailOnUnresolved bool
}

// DefaultOptions returns sensible defaults for Options.
func DefaultOptions() Options {
	return Options{
		MaxPasses:        10,
		FailOnUnresolved: false,
	}
}

// Interpolate expands all inter-variable references in env using the values
// within env itself (and falling back to OS env via expander.Expand).
// It mutates and returns the same map for convenience.
func Interpolate(env map[string]string, opts Options) (map[string]string, error) {
	if opts.MaxPasses <= 0 {
		opts.MaxPasses = DefaultOptions().MaxPasses
	}

	for pass := 0; pass < opts.MaxPasses; pass++ {
		changed := false

		for key, val := range env {
			if !expander.HasReference(val) {
				continue
			}

			expanded := expander.Expand(val, env)
			if expanded != val {
				env[key] = expanded
				changed = true
			}
		}

		if !changed {
			break
		}
	}

	if opts.FailOnUnresolved {
		var unresolved []string
		for key, val := range env {
			if expander.HasReference(val) {
				unresolved = append(unresolved, key)
			}
		}
		if len(unresolved) > 0 {
			return env, fmt.Errorf("unresolved references in keys: %s", strings.Join(unresolved, ", "))
		}
	}

	return env, nil
}
