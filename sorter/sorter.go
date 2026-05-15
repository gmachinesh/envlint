// Package sorter provides utilities for sorting and ordering
// environment variable keys and validation results.
package sorter

import (
	"sort"
	"strings"

	"github.com/nicholasgasior/envlint/validator"
)

// Order defines how keys should be sorted.
type Order string

const (
	// OrderAsc sorts keys in ascending (A-Z) order.
	OrderAsc Order = "asc"
	// OrderDesc sorts keys in descending (Z-A) order.
	OrderDesc Order = "desc"
	// OrderSeverity sorts results by severity: error > warning > info > ok.
	OrderSeverity Order = "severity"
)

// Options configures sorting behaviour.
type Options struct {
	Order Order
}

// DefaultOptions returns sensible sorting defaults.
func DefaultOptions() Options {
	return Options{Order: OrderAsc}
}

// SortKeys returns a sorted copy of the provided key slice.
func SortKeys(keys []string, opts Options) []string {
	copy_ := make([]string, len(keys))
	copy(copy_, keys)

	switch opts.Order {
	case OrderDesc:
		sort.Slice(copy_, func(i, j int) bool {
			return strings.ToLower(copy_[i]) > strings.ToLower(copy_[j])
		})
	default:
		sort.Slice(copy_, func(i, j int) bool {
			return strings.ToLower(copy_[i]) < strings.ToLower(copy_[j])
		})
	}

	return copy_
}

// severityRank maps severity strings to numeric rank (lower = higher priority).
func severityRank(s string) int {
	switch strings.ToLower(s) {
	case "error":
		return 0
	case "warning":
		return 1
	case "info":
		return 2
	default:
		return 3
	}
}

// SortResults returns a sorted copy of the validation results slice.
// When OrderSeverity is used, results are sorted by severity then key.
// Otherwise they are sorted alphabetically by key.
func SortResults(results []validator.Result, opts Options) []validator.Result {
	copy_ := make([]validator.Result, len(results))
	copy(copy_, results)

	switch opts.Order {
	case OrderSeverity:
		sort.SliceStable(copy_, func(i, j int) bool {
			ri, rj := severityRank(copy_[i].Severity), severityRank(copy_[j].Severity)
			if ri != rj {
				return ri < rj
			}
			return strings.ToLower(copy_[i].Key) < strings.ToLower(copy_[j].Key)
		})
	case OrderDesc:
		sort.SliceStable(copy_, func(i, j int) bool {
			return strings.ToLower(copy_[i].Key) > strings.ToLower(copy_[j].Key)
		})
	default:
		sort.SliceStable(copy_, func(i, j int) bool {
			return strings.ToLower(copy_[i].Key) < strings.ToLower(copy_[j].Key)
		})
	}

	return copy_
}
