package auditor_test

import (
	"testing"

	"github.com/user/envlint/auditor"
	"github.com/user/envlint/schema"
)

func makeSchema(names ...string) *schema.Schema {
	s := &schema.Schema{}
	for _, n := range names {
		s.Vars = append(s.Vars, schema.Var{Name: n})
	}
	return s
}

func TestAudit_NoUnknown(t *testing.T) {
	env := map[string]string{"APP_PORT": "8080", "APP_HOST": "localhost"}
	s := makeSchema("APP_PORT", "APP_HOST")
	results := auditor.Audit(env, s)
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestAudit_OneUnknown(t *testing.T) {
	env := map[string]string{"APP_PORT": "8080", "UNKNOWN_KEY": "value"}
	s := makeSchema("APP_PORT")
	results := auditor.Audit(env, s)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "UNKNOWN_KEY" {
		t.Errorf("expected key UNKNOWN_KEY, got %q", results[0].Key)
	}
	if results[0].Severity != "warning" {
		t.Errorf("expected severity warning, got %q", results[0].Severity)
	}
}

func TestAudit_AllUnknown(t *testing.T) {
	env := map[string]string{"FOO": "1", "BAR": "2"}
	s := makeSchema()
	results := auditor.Audit(env, s)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestAudit_EmptyEnv(t *testing.T) {
	env := map[string]string{}
	s := makeSchema("APP_PORT")
	results := auditor.Audit(env, s)
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestHasUnknown_True(t *testing.T) {
	results := []auditor.Result{{Key: "X", Message: "msg", Severity: "warning"}}
	if !auditor.HasUnknown(results) {
		t.Error("expected HasUnknown to return true")
	}
}

func TestHasUnknown_False(t *testing.T) {
	if auditor.HasUnknown(nil) {
		t.Error("expected HasUnknown to return false")
	}
}
