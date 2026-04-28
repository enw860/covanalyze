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
