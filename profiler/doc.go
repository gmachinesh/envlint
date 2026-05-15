// Package profiler provides statistical analysis of .env variable maps.
//
// It categorises values by type (numeric, boolean, URL, empty) and groups
// keys by their underscore-delimited prefix, giving a quick overview of the
// structure and health of an environment configuration.
//
// Basic usage:
//
//	env := map[string]string{
//		"APP_PORT":  "8080",
//		"APP_DEBUG": "true",
//		"DB_HOST":   "localhost",
//	}
//
//	p := profiler.Analyze(env)
//	profiler.ReportText(p, os.Stdout)
//
// The Profile struct is also JSON-serialisable via ReportJSON for
// machine-readable output in CI pipelines.
package profiler
