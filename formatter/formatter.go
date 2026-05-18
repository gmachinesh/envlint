package formatter

import (
	"fmt"
	"strings"
)

// Style represents the output format for variable names and messages.
type Style int

const (
	StyleDefault Style = iota
	StyleUppercase
	StyleLowercase
	StyleSnakeCase
)

// Options configures how messages and variable names are formatted.
type Options struct {
	Style       Style
	Prefix      string
	MaskValues  bool
	MaxValueLen int
}

// DefaultOptions returns sensible formatting defaults.
func DefaultOptions() Options {
	return Options{
		Style:       StyleDefault,
		MaskValues:  false,
		MaxValueLen: 64,
	}
}

// FormatKey applies the configured style to an environment variable key.
func FormatKey(key string, opts Options) string {
	switch opts.Style {
	case StyleUppercase:
		key = strings.ToUpper(key)
	case StyleLowercase:
		key = strings.ToLower(key)
	case StyleSnakeCase:
		key = toSnakeCase(key)
	}
	if opts.Prefix != "" {
		return fmt.Sprintf("%s%s", opts.Prefix, key)
	}
	return key
}

// FormatValue returns a display-safe version of a value.
func FormatValue(value string, opts Options) string {
	if opts.MaskValues {
		if len(value) == 0 {
			return ""
		}
		return "***"
	}
	if opts.MaxValueLen > 0 && len(value) > opts.MaxValueLen {
		return value[:opts.MaxValueLen] + "..."
	}
	return value
}

// FormatKeyValue formats a key-value pair as a single string using the
// provided options for both the key and the value. The result takes the
// form "KEY=value", consistent with standard environment variable notation.
func FormatKeyValue(key, value string, opts Options) string {
	return fmt.Sprintf("%s=%s", FormatKey(key, opts), FormatValue(value, opts))
}

// toSnakeCase converts a string to snake_case.
func toSnakeCase(s string) string {
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, " ", "_")
	return strings.ToLower(s)
}
