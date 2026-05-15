// Package scorer computes a numeric health score (0–100) for a .env file
// based on the results produced by the validator package.
//
// # Scoring model
//
// Each validator result carries a severity level. The scorer deducts a
// configurable penalty for every result and clamps the final value to the
// [0, 100] range:
//
//	score = max(0, 100 - Σ(count(severity) × weight(severity)))
//
// Default weights:
//
//	error   → −10 points
//	warning → −3  points
//	info    → −1  point
//
// # Letter grades
//
// The Grade helper maps a numeric score to a familiar letter grade:
//
//	90–100 → A
//	75–89  → B
//	60–74  → C
//	40–59  → D
//	0–39   → F
package scorer
