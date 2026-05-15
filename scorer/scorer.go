// Package scorer assigns a numeric health score to a .env file based on
// validator results. A perfect score of 100 means no issues were found;
// each error, warning, and info result deducts a weighted penalty.
package scorer

import (
	"math"

	"github.com/user/envlint/validator"
)

// DefaultWeights defines the penalty deducted per result severity.
var DefaultWeights = Weights{
	Error:   10,
	Warning: 3,
	Info:    1,
}

// Weights controls how many points each severity level deducts.
type Weights struct {
	Error   int
	Warning int
	Info    int
}

// Score holds the computed health score together with a breakdown.
type Score struct {
	// Value is the final score in the range [0, 100].
	Value int
	// Penalties is the total points deducted before clamping.
	Penalties int
	ErrorCount   int
	WarningCount int
	InfoCount    int
}

// Compute calculates a health score from the supplied validator results.
// An empty result slice returns a perfect score of 100.
func Compute(results []validator.Result, w Weights) Score {
	var errs, warns, infos int
	for _, r := range results {
		switch r.Severity {
		case "error":
			errs++
		case "warning":
			warns++
		case "info":
			infos++
		}
	}

	penalty := errs*w.Error + warns*w.Warning + infos*w.Info
	value := int(math.Max(0, float64(100-penalty)))

	return Score{
		Value:        value,
		Penalties:    penalty,
		ErrorCount:   errs,
		WarningCount: warns,
		InfoCount:    infos,
	}
}

// Grade returns a letter grade for a given score value.
func Grade(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 75:
		return "B"
	case score >= 60:
		return "C"
	case score >= 40:
		return "D"
	default:
		return "F"
	}
}
