// Package resolver resolves the final value of an environment variable
// by walking a priority-ordered list of sources (env map, OS environment,
// and schema-defined defaults).
package resolver

import (
	"fmt"
	"os"

	"github.com/user/envlint/schema"
)

// Options controls resolver behaviour.
type Options struct {
	// FallbackToOS allows falling back to os.Getenv when a key is absent
	// from all provided env maps.
	FallbackToOS bool
	// FallbackToDefault uses the schema default when all other sources miss.
	FallbackToDefault bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		FallbackToOS:      true,
		FallbackToDefault: true,
	}
}

// Result holds the resolved value and the source it came from.
type Result struct {
	Key    string
	Value  string
	Source string // "env", "os", "default", or "missing"
	Found  bool
}

// Resolve returns the resolved Result for key given the priority:
// env map → OS environment (optional) → schema default (optional).
func Resolve(key string, env map[string]string, s *schema.Schema, opts Options) Result {
	if v, ok := env[key]; ok {
		return Result{Key: key, Value: v, Source: "env", Found: true}
	}

	if opts.FallbackToOS {
		if v, ok := os.LookupEnv(key); ok {
			return Result{Key: key, Value: v, Source: "os", Found: true}
		}
	}

	if opts.FallbackToDefault && s != nil {
		for _, entry := range s.Vars {
			if entry.Name == key && entry.Default != "" {
				return Result{Key: key, Value: entry.Default, Source: "default", Found: true}
			}
		}
	}

	return Result{Key: key, Source: "missing", Found: false}
}

// ResolveAll resolves every key defined in the schema.
func ResolveAll(env map[string]string, s *schema.Schema, opts Options) []Result {
	if s == nil {
		return nil
	}
	results := make([]Result, 0, len(s.Vars))
	for _, entry := range s.Vars {
		results = append(results, Resolve(entry.Name, env, s, opts))
	}
	return results
}

// MustResolve resolves a key and returns an error if the key is not found.
func MustResolve(key string, env map[string]string, s *schema.Schema, opts Options) (Result, error) {
	r := Resolve(key, env, s, opts)
	if !r.Found {
		return r, fmt.Errorf("resolver: key %q could not be resolved from any source", key)
	}
	return r, nil
}
