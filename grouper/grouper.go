// Package grouper organises validation results by a chosen dimension
// (prefix, severity, or status) to make large reports easier to scan.
package grouper

import (
	"sort"
	"strings"

	"github.com/yourorg/envlint/validator"
)

// By is the dimension used to group results.
type By string

const (
	ByPrefix   By = "prefix"
	BySeverity By = "severity"
	ByStatus   By = "status"
)

// Group holds a named collection of results.
type Group struct {
	Name    string
	Results []validator.Result
}

// Apply partitions results into named groups according to the chosen dimension.
// Groups are returned in deterministic (sorted) order.
func Apply(results []validator.Result, by By) []Group {
	buckets := map[string][]validator.Result{}

	for _, r := range results {
		key := bucketKey(r, by)
		buckets[key] = append(buckets[key], r)
	}

	groups := make([]Group, 0, len(buckets))
	for name, rs := range buckets {
		groups = append(groups, Group{Name: name, Results: rs})
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})

	return groups
}

func bucketKey(r validator.Result, by By) string {
	switch by {
	case BySeverity:
		return r.Severity
	case ByStatus:
		if r.OK {
			return "ok"
		}
		return "fail"
	case ByPrefix:
		parts := strings.SplitN(r.Key, "_", 2)
		if len(parts) > 1 {
			return strings.ToLower(parts[0])
		}
		return "(none)"
	default:
		return "(unknown)"
	}
}
