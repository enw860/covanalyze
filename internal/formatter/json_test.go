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
