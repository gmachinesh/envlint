package resolver_test

import (
	"os"
	"testing"

	"github.com/user/envlint/resolver"
	"github.com/user/envlint/schema"
)

func makeSchema() *schema.Schema {
	return &schema.Schema{
		Vars: []schema.VarDef{
			{Name: "APP_HOST", Required: true, Default: "localhost"},
			{Name: "APP_PORT", Required: true, Default: ""},
			{Name: "APP_SECRET", Required: true, Default: ""},
		},
	}
}

func TestResolve_FromEnvMap(t *testing.T) {
	env := map[string]string{"APP_HOST": "example.com"}
	r := resolver.Resolve("APP_HOST", env, makeSchema(), resolver.DefaultOptions())
	if !r.Found || r.Value != "example.com" || r.Source != "env" {
		t.Fatalf("expected env source, got %+v", r)
	}
}

func TestResolve_FallbackToOS(t *testing.T) {
	os.Setenv("APP_PORT", "9090")
	defer os.Unsetenv("APP_PORT")

	env := map[string]string{}
	opts := resolver.DefaultOptions()
	r := resolver.Resolve("APP_PORT", env, makeSchema(), opts)
	if !r.Found || r.Value != "9090" || r.Source != "os" {
		t.Fatalf("expected os source, got %+v", r)
	}
}

func TestResolve_FallbackToDefault(t *testing.T) {
	env := map[string]string{}
	opts := resolver.Options{FallbackToOS: false, FallbackToDefault: true}
	r := resolver.Resolve("APP_HOST", env, makeSchema(), opts)
	if !r.Found || r.Value != "localhost" || r.Source != "default" {
		t.Fatalf("expected default source, got %+v", r)
	}
}

func TestResolve_Missing(t *testing.T) {
	env := map[string]string{}
	opts := resolver.Options{FallbackToOS: false, FallbackToDefault: false}
	r := resolver.Resolve("APP_SECRET", env, makeSchema(), opts)
	if r.Found {
		t.Fatalf("expected not found, got %+v", r)
	}
	if r.Source != "missing" {
		t.Fatalf("expected source=missing, got %q", r.Source)
	}
}

func TestResolveAll_ReturnsAllKeys(t *testing.T) {
	env := map[string]string{"APP_PORT": "8080"}
	opts := resolver.Options{FallbackToOS: false, FallbackToDefault: true}
	results := resolver.ResolveAll(env, makeSchema(), opts)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestMustResolve_Found(t *testing.T) {
	env := map[string]string{"APP_HOST": "prod.example.com"}
	r, err := resolver.MustResolve("APP_HOST", env, makeSchema(), resolver.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Value != "prod.example.com" {
		t.Fatalf("unexpected value: %q", r.Value)
	}
}

func TestMustResolve_NotFound(t *testing.T) {
	env := map[string]string{}
	opts := resolver.Options{FallbackToOS: false, FallbackToDefault: false}
	_, err := resolver.MustResolve("APP_SECRET", env, makeSchema(), opts)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
