package formatter

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/enw860/covanalyze/internal/models"
)

func TestFormatJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     *models.Output
		expectErr bool
	}{
		{
			name:      "nil output returns validation error",
			input:     nil,
			expectErr: true,
		},
		{
			name: "successful JSON formatting with single file",
			input: &models.Output{
				FileReports: []models.FileReport{
					{
						File:         "main.go",
						CoverageRate: 0.8500,
						UncoveredItems: []models.UncoveredItem{
							{LineRange: "10-15"},
						},
					},
				},
			},
			expectErr: false,
		},
		{
			name: "empty file_reports array",
			input: &models.Output{
				FileReports: []models.FileReport{},
			},
			expectErr: false,
		},
		{
			name: "files with 100% coverage (empty uncovered_items)",
			input: &models.Output{
				FileReports: []models.FileReport{
					{
						File:           "perfect.go",
						CoverageRate:   1.0000,
						UncoveredItems: []models.UncoveredItem{},
					},
				},
			},
			expectErr: false,
		},
		{
			name: "multiple files with various coverage rates",
			input: &models.Output{
				FileReports: []models.FileReport{
					{
						File:         "partial.go",
						CoverageRate: 0.7500,
						UncoveredItems: []models.UncoveredItem{
							{LineRange: "10"},
							{LineRange: "20-25"},
						},
					},
					{
						File:           "full.go",
						CoverageRate:   1.0000,
						UncoveredItems: []models.UncoveredItem{},
					},
					{
						File:         "none.go",
						CoverageRate: 0.0000,
						UncoveredItems: []models.UncoveredItem{
							{LineRange: "1-100"},
						},
					},
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := FormatJSON(tt.input)

			// Check error expectation
			if tt.expectErr {
				if err == nil {
					t.Error("FormatJSON() expected error, got nil")
				}
				return
			}

			// No error expected
			if err != nil {
				t.Errorf("FormatJSON() unexpected error: %v", err)
				return
			}

			// Verify it's valid JSON
			var output models.Output
			if err := json.Unmarshal(jsonBytes, &output); err != nil {
				t.Errorf("Invalid JSON: %v", err)
			}
		})
	}
}

func TestFormatJSON_SchemaReference(t *testing.T) {
	input := &models.Output{
		FileReports: []models.FileReport{},
	}

	jsonBytes, err := FormatJSON(input)
	if err != nil {
		t.Fatalf("FormatJSON() error: %v", err)
	}

	var output models.Output
	if err := json.Unmarshal(jsonBytes, &output); err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}

	if output.Schema != "./schema/output.schema.json" {
		t.Errorf("Schema = %q, want %q", output.Schema, "./schema/output.schema.json")
	}
}

func TestFormatJSON_Sorting(t *testing.T) {
	input := &models.Output{
		FileReports: []models.FileReport{
			{File: "low.go", CoverageRate: 0.3000, UncoveredItems: []models.UncoveredItem{}},
			{File: "high.go", CoverageRate: 0.9500, UncoveredItems: []models.UncoveredItem{}},
			{File: "medium.go", CoverageRate: 0.6000, UncoveredItems: []models.UncoveredItem{}},
		},
	}

	jsonBytes, err := FormatJSON(input)
	if err != nil {
		t.Fatalf("FormatJSON() error: %v", err)
	}

	var output models.Output
	if err := json.Unmarshal(jsonBytes, &output); err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}

	// Verify sorting: high.go (0.95), medium.go (0.60), low.go (0.30)
	if len(output.FileReports) != 3 {
		t.Fatalf("FileReports length = %d, want 3", len(output.FileReports))
	}
	if output.FileReports[0].File != "high.go" {
		t.Errorf("FileReports[0].File = %q, want %q", output.FileReports[0].File, "high.go")
	}
	if output.FileReports[1].File != "medium.go" {
		t.Errorf("FileReports[1].File = %q, want %q", output.FileReports[1].File, "medium.go")
	}
	if output.FileReports[2].File != "low.go" {
		t.Errorf("FileReports[2].File = %q, want %q", output.FileReports[2].File, "low.go")
	}
}

func TestFormatJSON_Indentation(t *testing.T) {
	input := &models.Output{
		FileReports: []models.FileReport{
			{
				File:         "test.go",
				CoverageRate: 0.5000,
				UncoveredItems: []models.UncoveredItem{
					{LineRange: "10-15"},
				},
			},
		},
	}

	jsonBytes, err := FormatJSON(input)
	if err != nil {
		t.Fatalf("FormatJSON() error: %v", err)
	}

	jsonStr := string(jsonBytes)

	// Check for 2-space indentation at root level
	if !strings.Contains(jsonStr, "  \"$schema\"") {
		t.Error("JSON not indented with 2 spaces at root level")
	}
	if !strings.Contains(jsonStr, "  \"file_reports\"") {
		t.Error("JSON not indented with 2 spaces at root level")
	}

	// Check for nested indentation (4 spaces for nested objects)
	if !strings.Contains(jsonStr, "    \"file\"") {
		t.Error("Nested JSON not indented with 4 spaces")
	}
	if !strings.Contains(jsonStr, "    \"coverage_rate\"") {
		t.Error("Nested JSON not indented with 4 spaces")
	}
}

func TestFormatJSON_EmptyArrays(t *testing.T) {
	input := &models.Output{
		FileReports: []models.FileReport{
			{
				File:           "perfect.go",
				CoverageRate:   1.0000,
				UncoveredItems: []models.UncoveredItem{},
			},
		},
	}

	jsonBytes, err := FormatJSON(input)
	if err != nil {
		t.Fatalf("FormatJSON() error: %v", err)
	}

	jsonStr := string(jsonBytes)

	// Verify empty arrays are rendered as [], not null
	if !strings.Contains(jsonStr, `"uncovered_items": []`) {
		t.Error("uncovered_items should be empty array [], not null")
	}
}

func TestFormatJSON_EnrichedUncoveredItems(t *testing.T) {
	tests := []struct {
		name  string
		input *models.Output
	}{
		{
			name: "enriched item with all fields populated",
			input: &models.Output{
				FileReports: []models.FileReport{
					{
						File:         "enriched.go",
						CoverageRate: 0.7500,
						UncoveredItems: []models.UncoveredItem{
							{
								Function:  "processData",
								Type:      "branch",
								Condition: "err != nil",
								LineRange: "10-15",
							},
						},
					},
				},
			},
		},
		{
			name: "enriched item with branch type",
			input: &models.Output{
				FileReports: []models.FileReport{
					{
						File:         "branch.go",
						CoverageRate: 0.8000,
						UncoveredItems: []models.UncoveredItem{
							{
								Function:  "handleRequest",
								Type:      "branch",
								Condition: "status == 200",
								LineRange: "20",
							},
						},
					},
				},
			},
		},
		{
			name: "enriched item with error_path type",
			input: &models.Output{
				FileReports: []models.FileReport{
					{
						File:         "error.go",
						CoverageRate: 0.6500,
						UncoveredItems: []models.UncoveredItem{
							{
								Function:  "validateInput",
								Type:      "error_path",
								Condition: "err != nil",
								LineRange: "30-35",
							},
						},
					},
				},
			},
		},
		{
			name: "enriched item with early_return type",
			input: &models.Output{
				FileReports: []models.FileReport{
					{
						File:         "return.go",
						CoverageRate: 0.9000,
						UncoveredItems: []models.UncoveredItem{
							{
								Function:  "checkCondition",
								Type:      "early_return",
								Condition: "value < 0",
								LineRange: "40",
							},
						},
					},
				},
			},
		},
		{
			name: "enriched item with loop type",
			input: &models.Output{
				FileReports: []models.FileReport{
					{
						File:         "loop.go",
						CoverageRate: 0.7000,
						UncoveredItems: []models.UncoveredItem{
							{
								Function:  "iterateItems",
								Type:      "loop",
								Condition: "i < len(items)",
								LineRange: "50-60",
							},
						},
					},
				},
			},
		},
		{
			name: "enriched item with empty type",
			input: &models.Output{
				FileReports: []models.FileReport{
					{
						File:         "unknown.go",
						CoverageRate: 0.5000,
						UncoveredItems: []models.UncoveredItem{
							{
								Function:  "unknownContext",
								Type:      "",
								Condition: "",
								LineRange: "70",
							},
						},
					},
				},
			},
		},
		{
			name: "multiple enriched items with different types",
			input: &models.Output{
				FileReports: []models.FileReport{
					{
						File:         "mixed.go",
						CoverageRate: 0.4500,
						UncoveredItems: []models.UncoveredItem{
							{
								Function:  "funcA",
								Type:      "branch",
								Condition: "x > 0",
								LineRange: "10",
							},
							{
								Function:  "funcB",
								Type:      "error_path",
								Condition: "err != nil",
								LineRange: "20-25",
							},
							{
								Function:  "funcC",
								Type:      "early_return",
								Condition: "done",
								LineRange: "30",
							},
							{
								Function:  "funcD",
								Type:      "loop",
								Condition: "range items",
								LineRange: "40-50",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := FormatJSON(tt.input)
			if err != nil {
				t.Fatalf("FormatJSON() error: %v", err)
			}

			// Verify it's valid JSON
			var output models.Output
			if err := json.Unmarshal(jsonBytes, &output); err != nil {
				t.Fatalf("Invalid JSON: %v", err)
			}

			// Verify enriched fields are preserved
			if len(output.FileReports) != len(tt.input.FileReports) {
				t.Fatalf("FileReports length = %d, want %d", len(output.FileReports), len(tt.input.FileReports))
			}

			for i, fr := range output.FileReports {
				if len(fr.UncoveredItems) != len(tt.input.FileReports[i].UncoveredItems) {
					t.Fatalf("UncoveredItems length = %d, want %d", len(fr.UncoveredItems), len(tt.input.FileReports[i].UncoveredItems))
				}

				for j, item := range fr.UncoveredItems {
					expected := tt.input.FileReports[i].UncoveredItems[j]
					if item.Function != expected.Function {
						t.Errorf("UncoveredItem[%d].Function = %q, want %q", j, item.Function, expected.Function)
					}
					if item.Type != expected.Type {
						t.Errorf("UncoveredItem[%d].Type = %q, want %q", j, item.Type, expected.Type)
					}
					if item.Condition != expected.Condition {
						t.Errorf("UncoveredItem[%d].Condition = %q, want %q", j, item.Condition, expected.Condition)
					}
					if item.LineRange != expected.LineRange {
						t.Errorf("UncoveredItem[%d].LineRange = %q, want %q", j, item.LineRange, expected.LineRange)
					}
				}
			}
		})
	}
}

func TestFormatJSON_FunctionTruncation(t *testing.T) {
	// Create a function name exactly at 256 characters
	longFuncName := strings.Repeat("a", 256)

	// Create a pre-truncated function name (as enricher would produce)
	truncatedFuncName := strings.Repeat("b", 253) + "..."

	tests := []struct {
		name         string
		functionName string
		description  string
	}{
		{
			name:         "function name at 256 character limit",
			functionName: longFuncName,
			description:  "formatter preserves function at max length",
		},
		{
			name:         "pre-truncated function name with ellipsis",
			functionName: truncatedFuncName,
			description:  "formatter preserves pre-truncated function from enricher",
		},
		{
			name:         "short function name unchanged",
			functionName: "shortFunc",
			description:  "formatter preserves short function name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &models.Output{
				FileReports: []models.FileReport{
					{
						File:         "test.go",
						CoverageRate: 0.5000,
						UncoveredItems: []models.UncoveredItem{
							{
								Function:  tt.functionName,
								Type:      "branch",
								Condition: "x > 0",
								LineRange: "10",
							},
						},
					},
				},
			}

			jsonBytes, err := FormatJSON(input)
			if err != nil {
				t.Fatalf("FormatJSON() error: %v", err)
			}

			var output models.Output
			if err := json.Unmarshal(jsonBytes, &output); err != nil {
				t.Fatalf("Invalid JSON: %v", err)
			}

			actualFunc := output.FileReports[0].UncoveredItems[0].Function
			if actualFunc != tt.functionName {
				t.Errorf("Function = %q (len=%d), want %q (len=%d)",
					actualFunc, len(actualFunc), tt.functionName, len(tt.functionName))
			}

			// Verify function length is within schema limits
			if len(actualFunc) > models.MaxFunctionNameLength {
				t.Errorf("Function length = %d exceeds max %d", len(actualFunc), models.MaxFunctionNameLength)
			}
		})
	}
}

func TestFormatJSON_ConditionTruncation(t *testing.T) {
	// Create a condition exactly at 512 characters
	longCondition := strings.Repeat("x", 512)

	// Create a pre-truncated condition (as enricher would produce)
	truncatedCondition := strings.Repeat("y", 509) + "..."

	tests := []struct {
		name        string
		condition   string
		description string
	}{
		{
			name:        "condition at 512 character limit",
			condition:   longCondition,
			description: "formatter preserves condition at max length",
		},
		{
			name:        "pre-truncated condition with ellipsis",
			condition:   truncatedCondition,
			description: "formatter preserves pre-truncated condition from enricher",
		},
		{
			name:        "short condition unchanged",
			condition:   "err != nil",
			description: "formatter preserves short condition",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &models.Output{
				FileReports: []models.FileReport{
					{
						File:         "test.go",
						CoverageRate: 0.5000,
						UncoveredItems: []models.UncoveredItem{
							{
								Function:  "testFunc",
								Type:      "branch",
								Condition: tt.condition,
								LineRange: "10",
							},
						},
					},
				},
			}

			jsonBytes, err := FormatJSON(input)
			if err != nil {
				t.Fatalf("FormatJSON() error: %v", err)
			}

			var output models.Output
			if err := json.Unmarshal(jsonBytes, &output); err != nil {
				t.Fatalf("Invalid JSON: %v", err)
			}

			actualCond := output.FileReports[0].UncoveredItems[0].Condition
			if actualCond != tt.condition {
				t.Errorf("Condition = %q (len=%d), want %q (len=%d)",
					actualCond, len(actualCond), tt.condition, len(tt.condition))
			}

			// Verify condition length is within schema limits
			if len(actualCond) > models.MaxConditionLength {
				t.Errorf("Condition length = %d exceeds max %d", len(actualCond), models.MaxConditionLength)
			}
		})
	}
}

func TestFormatJSON_TypeEnumValidation(t *testing.T) {
	validTypes := []string{"branch", "error_path", "early_return", "loop", ""}

	for _, validType := range validTypes {
		t.Run("valid_type_"+validType, func(t *testing.T) {
			input := &models.Output{
				FileReports: []models.FileReport{
					{
						File:         "test.go",
						CoverageRate: 0.5000,
						UncoveredItems: []models.UncoveredItem{
							{
								Function:  "testFunc",
								Type:      validType,
								Condition: "x > 0",
								LineRange: "10",
							},
						},
					},
				},
			}

			jsonBytes, err := FormatJSON(input)
			if err != nil {
				t.Fatalf("FormatJSON() error: %v", err)
			}

			var output models.Output
			if err := json.Unmarshal(jsonBytes, &output); err != nil {
				t.Fatalf("Invalid JSON: %v", err)
			}

			actualType := output.FileReports[0].UncoveredItems[0].Type
			if actualType != validType {
				t.Errorf("Type = %q, want %q", actualType, validType)
			}
		})
	}
}
