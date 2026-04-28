package models

// CoverageBlock represents a single coverage block from coverage.out file.
// All position values are 1-based to match the coverage.out format.
type CoverageBlock struct {
	File      string `json:"file"`
	StartLine int    `json:"start_line"`
	StartCol  int    `json:"start_col"`
	EndLine   int    `json:"end_line"`
	EndCol    int    `json:"end_col"`
	NumStmt   int    `json:"num_stmt"`
	Count     int    `json:"count"`
}
