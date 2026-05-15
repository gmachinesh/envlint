package differ_test

import (
	"testing"

	"github.com/yourorg/envlint/differ"
)

func TestCompare_NoDifferences(t *testing.T) {
	base := map[string]string{"FOO": "bar", "BAZ": "qux"}
	cand := map[string]string{"FOO": "bar", "BAZ": "qux"}

	diffs := differ.Compare(base, cand)
	if len(diffs) != 0 {
		t.Fatalf("expected no diffs, got %d", len(diffs))
	}
}

func TestCompare_AddedKey(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	cand := map[string]string{"FOO": "bar", "NEW": "val"}

	diffs := differ.Compare(base, cand)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Type != differ.Added || diffs[0].Key != "NEW" {
		t.Errorf("unexpected diff: %+v", diffs[0])
	}
	if diffs[0].NewVal != "val" {
		t.Errorf("expected NewVal=val, got %q", diffs[0].NewVal)
	}
}

func TestCompare_RemovedKey(t *testing.T) {
	base := map[string]string{"FOO": "bar", "OLD": "gone"}
	cand := map[string]string{"FOO": "bar"}

	diffs := differ.Compare(base, cand)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Type != differ.Removed || diffs[0].Key != "OLD" {
		t.Errorf("unexpected diff: %+v", diffs[0])
	}
	if diffs[0].OldVal != "gone" {
		t.Errorf("expected OldVal=gone, got %q", diffs[0].OldVal)
	}
}

func TestCompare_ChangedKey(t *testing.T) {
	base := map[string]string{"FOO": "old"}
	cand := map[string]string{"FOO": "new"}

	diffs := differ.Compare(base, cand)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	d := diffs[0]
	if d.Type != differ.Changed || d.Key != "FOO" {
		t.Errorf("unexpected diff: %+v", d)
	}
	if d.OldVal != "old" || d.NewVal != "new" {
		t.Errorf("expected old->new, got %q->%q", d.OldVal, d.NewVal)
	}
}

func TestHasChanges_True(t *testing.T) {
	base := map[string]string{"A": "1"}
	cand := map[string]string{"A": "2"}
	if !differ.HasChanges(base, cand) {
		t.Error("expected HasChanges to return true")
	}
}

func TestHasChanges_False(t *testing.T) {
	base := map[string]string{"A": "1"}
	cand := map[string]string{"A": "1"}
	if differ.HasChanges(base, cand) {
		t.Error("expected HasChanges to return false")
	}
}
