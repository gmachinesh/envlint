package renamer_test

import (
	"testing"

	"github.com/nicholasgasior/envlint/renamer"
)

func TestRename_NoOptions(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out := renamer.Rename(env, renamer.DefaultOptions())
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Fatalf("expected keys unchanged, got %v", out)
	}
}

func TestRename_StripPrefix(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080", "OTHER": "x"}
	opts := renamer.DefaultOptions()
	opts.StripPrefix = "APP_"
	out := renamer.Rename(env, opts)
	if out["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", out["HOST"])
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", out["PORT"])
	}
	if out["OTHER"] != "x" {
		t.Errorf("expected OTHER=x, got %q", out["OTHER"])
	}
}

func TestRename_AddPrefix(t *testing.T) {
	env := map[string]string{"HOST": "localhost"}
	opts := renamer.DefaultOptions()
	opts.AddPrefix = "SVC_"
	out := renamer.Rename(env, opts)
	if out["SVC_HOST"] != "localhost" {
		t.Errorf("expected SVC_HOST, got %v", out)
	}
}

func TestRename_ToUpper(t *testing.T) {
	env := map[string]string{"db_host": "localhost"}
	opts := renamer.DefaultOptions()
	opts.ToUpper = true
	out := renamer.Rename(env, opts)
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST, got %v", out)
	}
}

func TestRename_ToLower(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost"}
	opts := renamer.DefaultOptions()
	opts.ToLower = true
	out := renamer.Rename(env, opts)
	if out["db_host"] != "localhost" {
		t.Errorf("expected db_host, got %v", out)
	}
}

func TestRename_StripAndAddPrefix(t *testing.T) {
	env := map[string]string{"OLD_KEY": "val"}
	opts := renamer.Options{StripPrefix: "OLD_", AddPrefix: "NEW_"}
	out := renamer.Rename(env, opts)
	if out["NEW_KEY"] != "val" {
		t.Errorf("expected NEW_KEY=val, got %v", out)
	}
}

func TestRenameKey_Single(t *testing.T) {
	opts := renamer.Options{StripPrefix: "PRE_", ToUpper: true}
	got := renamer.RenameKey("pre_value", opts)
	// StripPrefix is case-sensitive; "pre_" != "PRE_" so key unchanged then uppercased
	if got != "PRE_VALUE" {
		t.Errorf("expected PRE_VALUE, got %q", got)
	}
}

func TestRename_DoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	opts := renamer.Options{AddPrefix: "X_"}
	renamer.Rename(env, opts)
	if _, ok := env["X_FOO"]; ok {
		t.Error("original map should not be mutated")
	}
}
