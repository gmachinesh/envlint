// Package renamer provides utilities for bulk-renaming environment variable
// keys according to a set of transformation rules.
package renamer

import (
	"strings"
)

// Options controls how keys are renamed.
type Options struct {
	// StripPrefix removes a leading prefix from every key, if present.
	StripPrefix string

	// AddPrefix prepends a string to every key.
	AddPrefix string

	// ToUpper forces all keys to upper-case after other transformations.
	ToUpper bool

	// ToLower forces all keys to lower-case after other transformations.
	ToLower bool
}

// DefaultOptions returns an Options with no transformations applied.
func DefaultOptions() Options {
	return Options{}
}

// Rename applies the given Options to every key in env and returns a new map.
// Values are copied unchanged. If two source keys would map to the same
// renamed key the last one (in iteration order) wins.
func Rename(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		newKey := renameKey(k, opts)
		out[newKey] = v
	}
	return out
}

// RenameKey applies the transformation rules to a single key.
func RenameKey(key string, opts Options) string {
	return renameKey(key, opts)
}

func renameKey(key string, opts Options) string {
	if opts.StripPrefix != "" {
		key = strings.TrimPrefix(key, opts.StripPrefix)
	}
	if opts.AddPrefix != "" {
		key = opts.AddPrefix + key
	}
	switch {
	case opts.ToUpper:
		key = strings.ToUpper(key)
	case opts.ToLower:
		key = strings.ToLower(key)
	}
	return key
}
