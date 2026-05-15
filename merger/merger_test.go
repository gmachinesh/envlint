package merger_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envlint/merger"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestMerge_SingleFile(t *testing.T) {
	path := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	env, err := merger.Merge([]string{path}, merger.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", env["FOO"])
	}
}

func TestMerge_LaterFileWins(t *testing.T) {
	p1 := writeTempEnv(t, "FOO=first\nONLY=base\n")
	p2 := writeTempEnv(t, "FOO=second\n")
	env, err := merger.Merge([]string{p1, p2}, merger.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["FOO"] != "second" {
		t.Errorf("expected FOO=second, got %q", env["FOO"])
	}
	if env["ONLY"] != "base" {
		t.Errorf("expected ONLY=base, got %q", env["ONLY"])
	}
}

func TestMerge_MissingFileFailOnMissing(t *testing.T) {
	opts := merger.DefaultOptions()
	_, err := merger.Merge([]string{"/nonexistent/.env"}, opts)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestMerge_MissingFileSkipped(t *testing.T) {
	p := writeTempEnv(t, "KEY=val\n")
	opts := merger.Options{FailOnMissing: false}
	env, err := merger.Merge([]string{filepath.Join(t.TempDir(), "missing.env"), p}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["KEY"] != "val" {
		t.Errorf("expected KEY=val, got %q", env["KEY"])
	}
}

func TestSources_TracksOrigin(t *testing.T) {
	p1 := writeTempEnv(t, "ALPHA=1\n")
	p2 := writeTempEnv(t, "BETA=2\nALPHA=overridden\n")
	srcs, err := merger.Sources([]string{p1, p2}, merger.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if srcs["ALPHA"] != p2 {
		t.Errorf("expected ALPHA sourced from p2, got %q", srcs["ALPHA"])
	}
	if srcs["BETA"] != p2 {
		t.Errorf("expected BETA sourced from p2, got %q", srcs["BETA"])
	}
}
