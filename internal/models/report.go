package models

// UncoveredItem represents a single uncovered code segment.
// Type field is a strict enum: "branch", "error_path", "early_return", "loop", or empty string.
// Function names are truncated at 256 chars with "..." suffix.
// Conditions are truncated at 512 chars with "..." suffix.
// LineRange accepts both "45" and "45-52" formats.
type UncoveredItem struct {
	Function  string `json:"function,omitempty"`
	Type      string `json:"type,omitempty"`
	Condition string `json:"condition,omitempty"`
	LineRange string `json:"line_range"`
}

// FileReport represents per-file coverage analysis results.
// CoverageRate is rounded to 4 decimal places (e.g., 0.7845).
// CoverageRate = 1.0 for files with 0 statements.
// File paths are output exactly as they appear in coverage.out (no normalization).
// UncoveredItems are sorted by line number ascending.
type FileReport struct {
	File           string          `json:"file"`
	CoverageRate   float64         `json:"coverage_rate"`
	UncoveredItems []UncoveredItem `json:"uncovered_items"`
}
