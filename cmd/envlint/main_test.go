package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return p
}

func TestMain_MissingSchemaFile(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "-schema", "nonexistent.yaml", "-env", ".env")
	cmd.Dir = "."
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected non-zero exit, got nil; output: %s", out)
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		if exitErr.ExitCode() != 2 {
			t.Errorf("expected exit code 2, got %d", exitErr.ExitCode())
		}
	}
}

func TestMain_ValidRun(t *testing.T) {
	dir := t.TempDir()
	writeTempFile(t, dir, "test.schema.yaml", `variables:
  - key: APP_ENV
    required: true
`)
	writeTempFile(t, dir, "test.env", "APP_ENV=production\n")

	cmd := exec.Command("go", "run", filepath.Join("..", "..", "cmd", "envlint"),
		"-schema", filepath.Join(dir, "test.schema.yaml"),
		"-env", filepath.Join(dir, "test.env"),
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected success, got error: %v\noutput: %s", err, out)
	}
}
