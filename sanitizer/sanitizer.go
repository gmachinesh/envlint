// Package sanitizer provides utilities for cleaning and sanitizing
// environment variable keys and values before further processing.
package sanitizer

import (
	"strings"
	"unicode"
)

// Options controls sanitizer behaviour.
type Options struct {
	// TrimSpace removes leading/trailing whitespace from keys and values.
	TrimSpace bool
	// RemoveControlChars strips non-printable control characters from values.
	RemoveControlChars bool
	// NormalizeNewlines replaces \r\n and bare \r with \n in values.
	NormalizeNewlines bool
	// CollapseInternalSpace collapses runs of internal whitespace in keys.
	CollapseInternalSpace bool
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		TrimSpace:             true,
		RemoveControlChars:    true,
		NormalizeNewlines:     true,
		CollapseInternalSpace: false,
	}
}

// Sanitize applies the given options to every key/value pair in env,
// returning a new map with cleaned entries. The original map is not mutated.
func Sanitize(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		cleanKey := sanitizeKey(k, opts)
		cleanVal := sanitizeValue(v, opts)
		if cleanKey != "" {
			out[cleanKey] = cleanVal
		}
	}
	return out
}

func sanitizeKey(k string, opts Options) string {
	if opts.TrimSpace {
		k = strings.TrimSpace(k)
	}
	if opts.CollapseInternalSpace {
		k = strings.Join(strings.Fields(k), "_")
	}
	return k
}

func sanitizeValue(v string, opts Options) string {
	if opts.NormalizeNewlines {
		v = strings.ReplaceAll(v, "\r\n", "\n")
		v = strings.ReplaceAll(v, "\r", "\n")
	}
	if opts.RemoveControlChars {
		v = strings.Map(func(r rune) rune {
			if unicode.IsControl(r) && r != '\n' && r != '\t' {
				return -1
			}
			return r
		}, v)
	}
	if opts.TrimSpace {
		v = strings.TrimSpace(v)
	}
	return v
}
