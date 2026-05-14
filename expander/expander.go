// Package expander resolves variable references within .env values.
// It supports ${VAR} and $VAR syntax, expanding references from the
// already-parsed env map or falling back to OS environment variables.
package expander

import (
	"os"
	"strings"
)

// Expand resolves variable references in all values of the provided env map.
// It performs a single-pass expansion; circular references are left unresolved.
func Expand(env map[string]string) map[string]string {
	result := make(map[string]string, len(env))
	for k, v := range env {
		result[k] = expandValue(v, env)
	}
	return result
}

// expandValue replaces ${VAR} and $VAR references in s using env, then os.Getenv.
func expandValue(s string, env map[string]string) string {
	return os.Expand(s, func(key string) string {
		if val, ok := env[key]; ok {
			return val
		}
		return os.Getenv(key)
	})
}

// HasReference reports whether the value contains a variable reference.
func HasReference(s string) bool {
	return strings.Contains(s, "$")
}
