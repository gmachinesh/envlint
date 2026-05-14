package linter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envlint/linter"
)

func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return p
}

const schemaYAML = `variables:
  - key: APP_ENV
    required: true
    pattern: "^(development|staging|production)$"
    severity: error
  - key: PORT
    required: false
    severity: warning
`

func TestRun_AllValid(t *testing.T) {
	schemaPath := writeTempFile(t, "schema.yaml", schemaYAML)
	envPath := writeTempFile(t, ".env", "APP_ENV=production\nPORT=8080\n")

	res, err := linter.Run(linter.Options{SchemaPath: schemaPath, EnvPath: envPath})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Summary.Errors != 0 {
		t.Errorf("expected 0 errors, got %d", res.Summary.Errors)
	}
}

func TestRun_MissingRequired(t *testing.T) {
	schemaPath := writeTempFile(t, "schema.yaml", schemaYAML)
	envPath := writeTempFile(t, ".env", "PORT=9090\n")

	res, err := linter.Run(linter.Options{SchemaPath: schemaPath, EnvPath: envPath})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Summary.Errors == 0 {
		t.Error("expected at least one error for missing required variable")
	}
}

func TestRun_BadSchemaPath(t *testing.T) {
	_, err := linter.Run(linter.Options{SchemaPath: "/no/such/schema.yaml", EnvPath: "/tmp/.env"})
	if err == nil {
		t.Error("expected error for missing schema file")
	}
}

func TestRun_ExpandVars(t *testing.T) {
	schemaYAMLExpand := `variables:
  - key: APP_ENV
    required: true
    pattern: "^production$"
    severity: error
`
	schemaPath := writeTempFile(t, "schema.yaml", schemaYAMLExpand)
	envPath := writeTempFile(t, ".env", "BASE=production\nAPP_ENV=${BASE}\n")

	res, err := linter.Run(linter.Options{SchemaPath: schemaPath, EnvPath: envPath, ExpandVars: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Summary.Errors != 0 {
		t.Errorf("expected 0 errors after expansion, got %d", res.Summary.Errors)
	}
}

func TestRun_SeverityFilter(t *testing.T) {
	schemaPath := writeTempFile(t, "schema.yaml", schemaYAML)
	envPath := writeTempFile(t, ".env", "APP_ENV=bad-value\n")

	res, err := linter.Run(linter.Options{
		SchemaPath: schemaPath,
		EnvPath:    envPath,
		Severity:   "error",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range res.Results {
		if r.Severity != "error" {
			t.Errorf("expected only error severity, got %q", r.Severity)
		}
	}
}
