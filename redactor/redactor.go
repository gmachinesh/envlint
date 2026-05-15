// Package redactor masks sensitive values in .env output to prevent
// accidental exposure of secrets in logs or CI output.
package redactor

import "strings"

// DefaultSensitivePatterns holds common key substrings that indicate
// a value should be redacted.
var DefaultSensitivePatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"CREDENTIALS",
	"AUTH",
}

const redactedPlaceholder = "[REDACTED]"

// Options configures redaction behaviour.
type Options struct {
	// Patterns is the list of key substrings (case-insensitive) that
	// trigger redaction. When nil, DefaultSensitivePatterns is used.
	Patterns []string
}

// Redact returns a copy of env where values whose keys match any
// sensitive pattern are replaced with [REDACTED].
func Redact(env map[string]string, opts Options) map[string]string {
	patterns := opts.Patterns
	if len(patterns) == 0 {
		patterns = DefaultSensitivePatterns
	}

	result := make(map[string]string, len(env))
	for k, v := range env {
		if isSensitive(k, patterns) {
			result[k] = redactedPlaceholder
		} else {
			result[k] = v
		}
	}
	return result
}

// IsSensitive reports whether key matches any of the provided patterns.
func IsSensitive(key string, patterns []string) bool {
	return isSensitive(key, patterns)
}

func isSensitive(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}
