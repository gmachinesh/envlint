// Package profiler analyses .env variable usage patterns and produces
// a statistical summary useful for understanding the shape of an environment.
package profiler

import (
	"sort"
	"strings"
)

// Profile holds aggregate statistics about a set of env variables.
type Profile struct {
	Total       int
	Empty       int
	Numeric     int
	Boolean     int
	URL         int
	LongestKey  string
	ShortestKey string
	Prefixes    map[string]int // top-level prefix (before first "_") -> count
}

// Analyze inspects the provided env map and returns a Profile.
func Analyze(env map[string]string) Profile {
	if len(env) == 0 {
		return Profile{Prefixes: map[string]int{}}
	}

	p := Profile{
		Prefixes: make(map[string]int),
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	p.Total = len(keys)
	p.LongestKey = keys[0]
	p.ShortestKey = keys[0]

	for _, k := range keys {
		v := env[k]

		if len(k) > len(p.LongestKey) {
			p.LongestKey = k
		}
		if len(k) < len(p.ShortestKey) {
			p.ShortestKey = k
		}

		if v == "" {
			p.Empty++
		}
		if isNumeric(v) {
			p.Numeric++
		}
		if isBoolean(v) {
			p.Boolean++
		}
		if isURL(v) {
			p.URL++
		}

		if idx := strings.Index(k, "_"); idx > 0 {
			prefix := k[:idx]
			p.Prefixes[prefix]++
		}
	}

	return p
}

func isNumeric(v string) bool {
	if v == "" {
		return false
	}
	for i, c := range v {
		if c == '-' && i == 0 {
			continue
		}
		if c == '.' {
			continue
		}
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func isBoolean(v string) bool {
	switch strings.ToLower(v) {
	case "true", "false", "yes", "no", "1", "0":
		return true
	}
	return false
}

func isURL(v string) bool {
	l := strings.ToLower(v)
	return strings.HasPrefix(l, "http://") || strings.HasPrefix(l, "https://")
}
