// Package linter provides a high-level Run function that orchestrates the full
// envlint pipeline:
//
//  1. Load the YAML schema from disk (schema.Load).
//  2. Load the .env file from disk (dotenv.Load).
//  3. Optionally expand variable references within values (expander.Expand).
//  4. Validate the resolved env map against the schema (validator.Validate).
//  5. Apply severity / prefix filters to the raw results (filter.Apply).
//  6. Return a consolidated Result containing the filtered findings and a
//     Summary of error / warning / info counts.
//
// Typical usage:
//
//	res, err := linter.Run(linter.Options{
//		SchemaPath: ".env.schema.yaml",
//		EnvPath:    ".env",
//		ExpandVars: true,
//		Severity:   "error",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	if res.Summary.Errors > 0 {
//		os.Exit(1)
//	}
package linter
