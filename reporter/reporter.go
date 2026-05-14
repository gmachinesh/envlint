package reporter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/envlint/validator"
)

// Format represents the output format for the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Reporter writes validation results to an output stream.
type Reporter struct {
	Writer io.Writer
	Format Format
}

// New creates a Reporter writing to stdout with the given format.
func New(format Format) *Reporter {
	return &Reporter{Writer: os.Stdout, Format: format}
}

// Report writes the validation results and returns true if there were errors.
func (r *Reporter) Report(results []validator.Result) bool {
	switch r.Format {
	case FormatJSON:
		return r.reportJSON(results)
	default:
		return r.reportText(results)
	}
}

func (r *Reporter) reportText(results []validator.Result) bool {
	hasErrors := false
	for _, res := range results {
		if res.Error != "" {
			hasErrors = true
			fmt.Fprintf(r.Writer, "[ERROR] %s: %s\n", res.Key, res.Error)
		} else {
			fmt.Fprintf(r.Writer, "[OK]    %s\n", res.Key)
		}
	}
	if len(results) == 0 {
		fmt.Fprintln(r.Writer, "No variables to validate.")
	}
	return hasErrors
}

func (r *Reporter) reportJSON(results []validator.Result) bool {
	hasErrors := false
	entries := make([]string, 0, len(results))
	for _, res := range results {
		status := "ok"
		errField := ""
		if res.Error != "" {
			status = "error"
			hasErrors = true
			errField = fmt.Sprintf(",\"error\":%q", res.Error)
		}
		entries = append(entries, fmt.Sprintf(`{"key":%q,"status":%q%s}`, res.Key, status, errField))
	}
	fmt.Fprintf(r.Writer, "[%s]\n", strings.Join(entries, ","))
	return hasErrors
}
