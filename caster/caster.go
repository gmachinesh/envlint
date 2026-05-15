// Package caster provides type-casting utilities for .env variable values.
// It converts raw string values into typed Go values such as int, bool, float64,
// and duration, returning descriptive errors on failure.
package caster

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Result holds the outcome of a cast attempt.
type Result struct {
	Key    string
	Raw    string
	Kind   string
	Value  any
	Err    error
}

// CastInt parses the raw string as a base-10 integer.
func CastInt(key, raw string) Result {
	v, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return Result{Key: key, Raw: raw, Kind: "int", Err: fmt.Errorf("%s: cannot cast %q to int", key, raw)}
	}
	return Result{Key: key, Raw: raw, Kind: "int", Value: v}
}

// CastBool parses the raw string as a boolean (1, t, true, 0, f, false).
func CastBool(key, raw string) Result {
	v, err := strconv.ParseBool(strings.TrimSpace(raw))
	if err != nil {
		return Result{Key: key, Raw: raw, Kind: "bool", Err: fmt.Errorf("%s: cannot cast %q to bool", key, raw)}
	}
	return Result{Key: key, Raw: raw, Kind: "bool", Value: v}
}

// CastFloat parses the raw string as a 64-bit float.
func CastFloat(key, raw string) Result {
	v, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		return Result{Key: key, Raw: raw, Kind: "float", Err: fmt.Errorf("%s: cannot cast %q to float", key, raw)}
	}
	return Result{Key: key, Raw: raw, Kind: "float", Value: v}
}

// CastDuration parses the raw string as a time.Duration (e.g. "5s", "2m").
func CastDuration(key, raw string) Result {
	v, err := time.ParseDuration(strings.TrimSpace(raw))
	if err != nil {
		return Result{Key: key, Raw: raw, Kind: "duration", Err: fmt.Errorf("%s: cannot cast %q to duration", key, raw)}
	}
	return Result{Key: key, Raw: raw, Kind: "duration", Value: v}
}

// Cast dispatches to the appropriate cast function based on kind.
// Supported kinds: "int", "bool", "float", "duration", "string".
func Cast(key, raw, kind string) Result {
	switch strings.ToLower(kind) {
	case "int":
		return CastInt(key, raw)
	case "bool":
		return CastBool(key, raw)
	case "float":
		return CastFloat(key, raw)
	case "duration":
		return CastDuration(key, raw)
	case "string", "":
		return Result{Key: key, Raw: raw, Kind: "string", Value: raw}
	default:
		return Result{Key: key, Raw: raw, Kind: kind, Err: fmt.Errorf("%s: unknown kind %q", key, kind)}
	}
}
