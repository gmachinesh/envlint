// Package coercer normalises raw .env string values into canonical forms
// before validation or casting — e.g. trimming whitespace, normalising
// booleans, and stripping surrounding quotes.
package coercer

import (
	"strings"
)

// Options controls which coercions are applied.
type Options struct {
	TrimSpace       bool // strip leading/trailing whitespace
	NormaliseBool   bool // convert 1/0/yes/no/on/off → true/false
	StripQuotes     bool // remove surrounding single or double quotes
	LowercaseBool   bool // lowercase the result when it looks like a boolean
}

// DefaultOptions returns the recommended set of coercions.
func DefaultOptions() Options {
	return Options{
		TrimSpace:     true,
		NormaliseBool: true,
		StripQuotes:   true,
		LowercaseBool: true,
	}
}

var boolTrue = map[string]struct{}{
	"1": {}, "yes": {}, "on": {}, "true": {},
}
var boolFalse = map[string]struct{}{
	"0": {}, "no": {}, "off": {}, "false": {},
}

// Coerce applies the enabled transformations to value and returns the result.
func Coerce(value string, opts Options) string {
	if opts.TrimSpace {
		value = strings.TrimSpace(value)
	}

	if opts.StripQuotes {
		value = stripQuotes(value)
	}

	if opts.NormaliseBool {
		lower := strings.ToLower(value)
		if _, ok := boolTrue[lower]; ok {
			if opts.LowercaseBool {
				return "true"
			}
			return "true"
		}
		if _, ok := boolFalse[lower]; ok {
			if opts.LowercaseBool {
				return "false"
			}
			return "false"
		}
	}

	return value
}

// CoerceAll applies Coerce to every value in env and returns a new map.
func CoerceAll(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = Coerce(v, opts)
	}
	return out
}

func stripQuotes(s string) string {
	if len(s) < 2 {
		return s
	}
	if (s[0] == '"' && s[len(s)-1] == '"') ||
		(s[0] == '\'' && s[len(s)-1] == '\'') {
		return s[1 : len(s)-1]
	}
	return s
}
