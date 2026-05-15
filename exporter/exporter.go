// Package exporter converts an env map into various output formats
// such as shell export statements, Docker --env-file syntax, and JSON.
package exporter

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Format represents the target output format.
type Format string

const (
	FormatShell  Format = "shell"
	FormatDocker Format = "docker"
	FormatJSON   Format = "json"
)

// Options controls exporter behaviour.
type Options struct {
	Format Format
	// SortKeys ensures deterministic output when true.
	SortKeys bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		Format:   FormatShell,
		SortKeys: true,
	}
}

// Export converts env into the requested format string.
// An error is returned only for the JSON format if marshalling fails.
func Export(env map[string]string, opts Options) (string, error) {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	if opts.SortKeys {
		sort.Strings(keys)
	}

	switch opts.Format {
	case FormatDocker:
		return exportDocker(keys, env), nil
	case FormatJSON:
		return exportJSON(keys, env)
	default:
		return exportShell(keys, env), nil
	}
}

func exportShell(keys []string, env map[string]string) string {
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "export %s=%q\n", k, env[k])
	}
	return sb.String()
}

func exportDocker(keys []string, env map[string]string) string {
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, env[k])
	}
	return sb.String()
}

func exportJSON(keys []string, env map[string]string) (string, error) {
	ordered := make(map[string]string, len(keys))
	for _, k := range keys {
		ordered[k] = env[k]
	}
	b, err := json.MarshalIndent(ordered, "", "  ")
	if err != nil {
		return "", fmt.Errorf("exporter: json marshal: %w", err)
	}
	return string(b) + "\n", nil
}
