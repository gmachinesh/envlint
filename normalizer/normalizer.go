// Package normalizer provides utilities for normalising environment variable
// keys and values into a consistent canonical form before validation or
// comparison.
package normalizer

import (
	"strings"
)

// Options controls how normalisation is applied.
type Options struct {
	// TrimSpace removes leading and trailing whitespace from keys and values.
	TrimSpace bool
	// UppercaseKeys converts all keys to UPPER_CASE.
	UppercaseKeys bool
	// LowercaseValues converts all values to lowercase.
	LowercaseValues bool
	// RemoveEmpty drops entries whose value is an empty string after trimming.
	RemoveEmpty bool
}

// DefaultOptions returns a sensible default normalisation configuration.
func DefaultOptions() Options {
	return Options{
		TrimSpace:       true,
		UppercaseKeys:   true,
		LowercaseValues: false,
		RemoveEmpty:     false,
	}
}

// Normalize applies the given Options to every entry in env and returns a new
// map containing the normalised key/value pairs. The original map is never
// mutated.
func Normalize(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if opts.TrimSpace {
			k = strings.TrimSpace(k)
			v = strings.TrimSpace(v)
		}
		if opts.UppercaseKeys {
			k = strings.ToUpper(k)
		}
		if opts.LowercaseValues {
			v = strings.ToLower(v)
		}
		if opts.RemoveEmpty && v == "" {
			continue
		}
		out[k] = v
	}
	return out
}

// NormalizeKey applies key-specific normalisation (trim + uppercase) to a
// single key string and returns the result.
func NormalizeKey(key string) string {
	return strings.ToUpper(strings.TrimSpace(key))
}
