// Package truncator provides utilities for truncating long environment
// variable values to a maximum length, useful for display and logging purposes.
package truncator

import "unicode/utf8"

// DefaultOptions returns a sensible set of truncation options.
func DefaultOptions() Options {
	return Options{
		MaxLen:  64,
		Suffix:  "...",
		OnlyValues: true,
	}
}

// Options controls how truncation is applied.
type Options struct {
	// MaxLen is the maximum number of runes allowed before truncation.
	MaxLen int

	// Suffix is appended when a value is truncated.
	Suffix string

	// OnlyValues, when true, skips truncation of keys.
	OnlyValues bool
}

// Truncate shortens s to at most opts.MaxLen runes, appending opts.Suffix
// when the string was actually shortened. If opts.MaxLen <= 0 the original
// string is returned unchanged.
func Truncate(s string, opts Options) string {
	if opts.MaxLen <= 0 {
		return s
	}
	if utf8.RuneCountInString(s) <= opts.MaxLen {
		return s
	}
	runes := []rune(s)
	return string(runes[:opts.MaxLen]) + opts.Suffix
}

// TruncateAll applies Truncate to every value (and optionally key) in env,
// returning a new map. The original map is never modified.
func TruncateAll(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		newKey := k
		if !opts.OnlyValues {
			newKey = Truncate(k, opts)
		}
		out[newKey] = Truncate(v, opts)
	}
	return out
}
