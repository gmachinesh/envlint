package caster

// BatchInput describes a single variable to cast in a batch operation.
type BatchInput struct {
	Key  string
	Raw  string
	Kind string
}

// BatchResult holds all results from a batch cast, separated into
// successes and failures for convenient downstream handling.
type BatchResult struct {
	Successes []Result
	Failures  []Result
}

// HasErrors returns true when at least one cast in the batch failed.
func (b BatchResult) HasErrors() bool {
	return len(b.Failures) > 0
}

// CastAll runs Cast over every input and partitions the results.
func CastAll(inputs []BatchInput) BatchResult {
	var br BatchResult
	for _, in := range inputs {
		r := Cast(in.Key, in.Raw, in.Kind)
		if r.Err != nil {
			br.Failures = append(br.Failures, r)
		} else {
			br.Successes = append(br.Successes, r)
		}
	}
	return br
}
