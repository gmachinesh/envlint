// Package caster provides type-casting utilities for raw .env string values.
//
// Supported target kinds:
//
//   - "string"  — identity cast; always succeeds
//   - "int"     — parsed as a base-10 integer via strconv.Atoi
//   - "bool"    — parsed via strconv.ParseBool (accepts 1/0, true/false, t/f)
//   - "float"   — parsed as a 64-bit IEEE-754 float
//   - "duration" — parsed as a time.Duration (e.g. "5s", "2m30s")
//
// Single-variable casting:
//
//	r := caster.Cast("PORT", "8080", "int")
//	if r.Err != nil {
//		log.Fatal(r.Err)
//	}
//	fmt.Println(r.Value.(int)) // 8080
//
// Batch casting:
//
//	br := caster.CastAll([]caster.BatchInput{
//		{Key: "PORT", Raw: "8080", Kind: "int"},
//		{Key: "DEBUG", Raw: "true", Kind: "bool"},
//	})
//	if br.HasErrors() {
//		for _, f := range br.Failures {
//			fmt.Println(f.Err)
//		}
//	}
package caster
