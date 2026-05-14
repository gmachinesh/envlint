package schema_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envlint/schema"
)

func TestLoad_ValidSchema(t *testing.T) {
	content := `
vars:
  DATABASE_URL:
    required: true
    type: url
    description: Primary database connection string
  DEBUG:
    required: false
    type: bool
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "schema.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	s, err := schema.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(s.Vars) != 2 {
		t.Errorf("expected 2 vars, got %d", len(s.Vars))
	}

	dbVar, ok := s.Vars["DATABASE_URL"]
	if !ok {
		t.Fatal("expected DATABASE_URL in schema")
	}
	if !dbVar.Required {
		t.Error("expected DATABASE_URL to be required")
	}
	if dbVar.Type != schema.TypeURL {
		t.Errorf("expected type url, got %s", dbVar.Type)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := schema.Load("/nonexistent/path/schema.yaml")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "bad.yaml")
	if err := os.WriteFile(path, []byte("{ invalid: yaml: ["), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := schema.Load(path)
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}

func TestLoad_EmptySchema(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "empty.yaml")
	if err := os.WriteFile(path, []byte("vars:\n"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	s, err := schema.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Vars == nil {
		t.Error("expected non-nil vars map")
	}
}
