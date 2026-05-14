package expander_test

import (
	"os"
	"testing"

	"github.com/user/envlint/expander"
)

func TestExpand_NoReferences(t *testing.T) {
	env := map[string]string{
		"FOO": "hello",
		"BAR": "world",
	}
	got := expander.Expand(env)
	if got["FOO"] != "hello" || got["BAR"] != "world" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestExpand_BraceReference(t *testing.T) {
	env := map[string]string{
		"BASE": "/app",
		"DATA": "${BASE}/data",
	}
	got := expander.Expand(env)
	if got["DATA"] != "/app/data" {
		t.Errorf("expected /app/data, got %q", got["DATA"])
	}
}

func TestExpand_UnbracedReference(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"URL":  "http://$HOST:8080",
	}
	got := expander.Expand(env)
	if got["URL"] != "http://localhost:8080" {
		t.Errorf("expected http://localhost:8080, got %q", got["URL"])
	}
}

func TestExpand_FallbackToOS(t *testing.T) {
	os.Setenv("OS_VAR", "from-os")
	defer os.Unsetenv("OS_VAR")

	env := map[string]string{
		"DERIVED": "${OS_VAR}-suffix",
	}
	got := expander.Expand(env)
	if got["DERIVED"] != "from-os-suffix" {
		t.Errorf("expected from-os-suffix, got %q", got["DERIVED"])
	}
}

func TestExpand_UnresolvedReference(t *testing.T) {
	env := map[string]string{
		"FOO": "${MISSING_VAR}",
	}
	got := expander.Expand(env)
	if got["FOO"] != "" {
		t.Errorf("expected empty string for unresolved ref, got %q", got["FOO"])
	}
}

func TestHasReference(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"plain", false},
		{"${VAR}", true},
		{"$VAR", true},
		{"", false},
	}
	for _, c := range cases {
		got := expander.HasReference(c.input)
		if got != c.want {
			t.Errorf("HasReference(%q) = %v, want %v", c.input, got, c.want)
		}
	}
}
