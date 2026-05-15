// Package merger provides utilities for combining multiple .env files into
// a single environment map.
//
// When the same key appears in more than one file, the value from the
// last file in the list takes precedence. This mirrors the behaviour of
// tools like docker-compose, which layer .env overrides on top of a base
// file.
//
// Basic usage:
//
//	env, err := merger.Merge([]string{".env", ".env.local"}, merger.DefaultOptions())
//	if err != nil {
//		log.Fatal(err)
//	}
//
// To understand which file last defined each key, use Sources:
//
//	srcs, err := merger.Sources([]string{".env", ".env.local"}, merger.DefaultOptions())
//	if err != nil {
//		log.Fatal(err)
//	}
//	for key, file := range srcs {
//		fmt.Printf("%s defined in %s\n", key, file)
//	}
package merger
