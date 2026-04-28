package models

// ASTNode represents a minimal AST node for enriching coverage data.
// Only essential fields are included - no column positions or parent references.
type ASTNode struct {
	NodeType      string `json:"node_type"`
	StartLine     int    `json:"start_line"`
	EndLine       int    `json:"end_line"`
	FunctionName  string `json:"function_name,omitempty"`
	ConditionText string `json:"condition_text,omitempty"`
}
