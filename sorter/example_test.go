package sorter_test

import (
	"fmt"

	"github.com/nicholasgasior/envlint/sorter"
	"github.com/nicholasgasior/envlint/validator"
)

func ExampleSortKeys() {
	keys := []string{"PORT", "APP_ENV", "DB_URL"}
	sorted := sorter.SortKeys(keys, sorter.DefaultOptions())
	for _, k := range sorted {
		fmt.Println(k)
	}
	// Output:
	// APP_ENV
	// DB_URL
	// PORT
}

func ExampleSortResults_severity() {
	results := []validator.Result{
		{Key: "LOG_LEVEL", Severity: "info"},
		{Key: "DB_PASS", Severity: "error"},
		{Key: "CACHE_TTL", Severity: "warning"},
	}
	sorted := sorter.SortResults(results, sorter.Options{Order: sorter.OrderSeverity})
	for _, r := range sorted {
		fmt.Printf("%s (%s)\n", r.Key, r.Severity)
	}
	// Output:
	// DB_PASS (error)
	// CACHE_TTL (warning)
	// LOG_LEVEL (info)
}
