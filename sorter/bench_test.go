package sorter_test

import (
	"fmt"
	"testing"

	"github.com/nicholasgasior/envlint/sorter"
	"github.com/nicholasgasior/envlint/validator"
)

func BenchmarkSortKeys_100(b *testing.B) {
	keys := make([]string, 100)
	for i := range keys {
		keys[i] = fmt.Sprintf("KEY_%03d", 100-i)
	}
	opts := sorter.DefaultOptions()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = sorter.SortKeys(keys, opts)
	}
}

func BenchmarkSortResults_Severity_100(b *testing.B) {
	severities := []string{"error", "warning", "info", "ok"}
	results := make([]validator.Result, 100)
	for i := range results {
		results[i] = validator.Result{
			Key:      fmt.Sprintf("KEY_%03d", 100-i),
			Severity: severities[i%4],
		}
	}
	opts := sorter.Options{Order: sorter.OrderSeverity}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = sorter.SortResults(results, opts)
	}
}
