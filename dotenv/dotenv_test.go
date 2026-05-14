package dotenv_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envlint/dotenv"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return path
}

func TestLoad_ValidEnv(t *testing.T) {
	path := writeTempEnv(t, `
# This is a comment
APP_ENV=production
DB_HOST=localhost
DB_PORT=5432
SECRET_KEY="my secret value"
ANOTHER='quoted'
`)
	env, err := dotenv.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cases := map[string]string{
		"APP_ENV":    "production",
		"DB_HOST":    "localhost",
		"DB_PORT":    "5432",
		"SECRET_KEY": "my secret value",
		"ANOTHER":    "quoted",
	}
	for k, want := range cases {
		if got := env[k]; got != want {
			t.Errorf("env[%q] = %q, want %q", k, got, want)
		}
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := dotenv.Load("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_MalformedLine(t *testing.T) {
	path := writeTempEnv(t, "VALID=ok\nBAD_LINE_NO_EQUALS\n")
	_, err := dotenv.Load(path)
	if err == nil {
		t.Fatal("expected error for malformed line, got nil")
	}
}

func TestLoad_InlineComment(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=myapp # the app name\n")
	env, err := dotenv.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := env["APP_NAME"]; got != "myapp" {
		t.Errorf("APP_NAME = %q, want %q", got, "myapp")
	}
}

func TestLoad_EmptyFile(t *testing.T) {
	path := writeTempEnv(t, "")
	env, err := dotenv.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 0 {
		t.Errorf("expected empty env, got %d entries", len(env))
	}
}
