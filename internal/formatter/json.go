package formatter

import (
	"encoding/json"
	"errors"
	"sort"

	"github.com/enw860/covanalyze/internal/models"
)

const (
	// SchemaReference is the hardcoded schema reference path
	SchemaReference = "./schema/output.schema.json"
	// JSONIndent is the indentation string for pretty-printed JSON
	JSONIndent = "  "
)

// FormatJSON formats the Output struct as pretty-printed JSON with validation.
// It sorts FileReports by coverage_rate descending and sets the schema reference.
// Returns the formatted JSON bytes and any validation error.
func FormatJSON(output *models.Output) ([]byte, error) {
	// Validate output is not nil
	if output == nil {
		return nil, errors.New("validation error: output cannot be nil")
	}

	// Set schema reference
	output.Schema = SchemaReference

	// Sort FileReports by coverage_rate descending
	sort.Slice(output.FileReports, func(i, j int) bool {
		return output.FileReports[i].CoverageRate > output.FileReports[j].CoverageRate
	})

	// Pretty-print JSON with 2-space indentation
	jsonBytes, err := json.MarshalIndent(output, "", JSONIndent)
	if err != nil {
		return nil, errors.New("validation error: failed to marshal JSON: " + err.Error())
	}

	return jsonBytes, nil
}
