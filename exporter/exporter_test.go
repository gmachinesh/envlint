package exporter_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/envlint/exporter"
)

var sampleEnv = map[string]string{
	"APP_ENV":  "production",
	"DB_HOST":  "localhost",
	"DB_PORT":  "5432",
}

func TestExport_Shell(t *testing.T) {
	out, err := exporter.Export(sampleEnv, exporter.Options{Format: exporter.FormatShell, SortKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `export APP_ENV="production"`) {
		t.Errorf("expected shell export line, got:\n%s", out)
	}
	if !strings.Contains(out, `export DB_HOST="localhost"`) {
		t.Errorf("missing DB_HOST export, got:\n%s", out)
	}
}

func TestExport_ShellSorted(t *testing.T) {
	out, _ := exporter.Export(sampleEnv, exporter.Options{Format: exporter.FormatShell, SortKeys: true})
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "export APP_ENV") {
		t.Errorf("expected APP_ENV first, got: %s", lines[0])
	}
}

func TestExport_Docker(t *testing.T) {
	out, err := exporter.Export(sampleEnv, exporter.Options{Format: exporter.FormatDocker, SortKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "export ") {
		t.Error("docker format should not contain 'export' keyword")
	}
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected KEY=VALUE line, got:\n%s", out)
	}
}

func TestExport_JSON(t *testing.T) {
	out, err := exporter.Export(sampleEnv, exporter.Options{Format: exporter.FormatJSON, SortKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]string
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if parsed["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %s", parsed["APP_ENV"])
	}
}

func TestExport_DefaultOptions(t *testing.T) {
	opts := exporter.DefaultOptions()
	if opts.Format != exporter.FormatShell {
		t.Errorf("expected default format shell, got %s", opts.Format)
	}
	if !opts.SortKeys {
		t.Error("expected SortKeys to be true by default")
	}
}

func TestExport_EmptyEnv(t *testing.T) {
	out, err := exporter.Export(map[string]string{}, exporter.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(out) != "" {
		t.Errorf("expected empty output, got: %q", out)
	}
}
