// Package exporter serialises an env map into common output formats.
//
// Supported formats:
//
//   - FormatShell  — POSIX shell export statements (export KEY="value")
//   - FormatDocker — Docker-compatible env-file syntax (KEY=value)
//   - FormatJSON   — Indented JSON object
//
// Basic usage:
//
//	env := map[string]string{"APP_ENV": "production", "PORT": "8080"}
//
//	// Shell export statements
//	out, err := exporter.Export(env, exporter.DefaultOptions())
//
//	// Docker env-file
//	out, err = exporter.Export(env, exporter.Options{
//		Format:   exporter.FormatDocker,
//		SortKeys: true,
//	})
//
//	// JSON
//	out, err = exporter.Export(env, exporter.Options{
//		Format:   exporter.FormatJSON,
//		SortKeys: true,
//	})
package exporter
