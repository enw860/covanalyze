package parser

import (
	"testing"

	"golang.org/x/tools/cover"
)

func TestCalculateFileCoverage(t *testing.T) {
	tests := []struct {
		name             string
		profile          *cover.Profile
		expectCoverage   float64
		expectUncovered  int
		expectLineRanges []string
	}{
		{
			name: "full coverage (100%)",
			profile: &cover.Profile{
				FileName: "example.go",
				Mode:     "set",
				Blocks: []cover.ProfileBlock{
					{StartLine: 10, StartCol: 2, EndLine: 10, EndCol: 20, NumStmt: 1, Count: 1},
					{StartLine: 15, StartCol: 2, EndLine: 17, EndCol: 3, NumStmt: 3, Count: 1},
					{StartLine: 20, StartCol: 2, EndLine: 20, EndCol: 15, NumStmt: 1, Count: 1},
				},
			},
			expectCoverage:   1.0,
			expectUncovered:  0,
			expectLineRanges: []string{},
		},
		{
			name: "partial coverage with single-line uncovered",
			profile: &cover.Profile{
				FileName: "partial.go",
				Mode:     "set",
				Blocks: []cover.ProfileBlock{
					{StartLine: 5, StartCol: 2, EndLine: 5, EndCol: 20, NumStmt: 1, Count: 1},
					{StartLine: 10, StartCol: 2, EndLine: 10, EndCol: 15, NumStmt: 1, Count: 0},
					{StartLine: 15, StartCol: 2, EndLine: 15, EndCol: 10, NumStmt: 1, Count: 1},
				},
			},
			expectCoverage:   0.6667,
			expectUncovered:  1,
			expectLineRanges: []string{"10"},
		},
		{
			name: "partial coverage with multi-line uncovered",
			profile: &cover.Profile{
				FileName: "multiline.go",
				Mode:     "set",
				Blocks: []cover.ProfileBlock{
					{StartLine: 5, StartCol: 2, EndLine: 5, EndCol: 20, NumStmt: 2, Count: 1},
					{StartLine: 10, StartCol: 2, EndLine: 15, EndCol: 3, NumStmt: 5, Count: 0},
					{StartLine: 20, StartCol: 2, EndLine: 20, EndCol: 10, NumStmt: 1, Count: 1},
				},
			},
			expectCoverage:   0.375,
			expectUncovered:  1,
			expectLineRanges: []string{"10-15"},
		},
		{
			name: "zero coverage (0%)",
			profile: &cover.Profile{
				FileName: "uncovered.go",
				Mode:     "set",
				Blocks: []cover.ProfileBlock{
					{StartLine: 5, StartCol: 2, EndLine: 5, EndCol: 20, NumStmt: 1, Count: 0},
					{StartLine: 10, StartCol: 2, EndLine: 12, EndCol: 3, NumStmt: 3, Count: 0},
					{StartLine: 15, StartCol: 2, EndLine: 15, EndCol: 10, NumStmt: 1, Count: 0},
				},
			},
			expectCoverage:   0.0,
			expectUncovered:  3,
			expectLineRanges: []string{"5", "10-12", "15"},
		},
		{
			name: "file with 0 statements",
			profile: &cover.Profile{
				FileName: "empty.go",
				Mode:     "set",
				Blocks:   []cover.ProfileBlock{},
			},
			expectCoverage:   1.0,
			expectUncovered:  0,
			expectLineRanges: []string{},
		},
		{
			name: "multiple uncovered blocks - test sorting",
			profile: &cover.Profile{
				FileName: "sorting.go",
				Mode:     "set",
				Blocks: []cover.ProfileBlock{
					{StartLine: 30, StartCol: 2, EndLine: 30, EndCol: 10, NumStmt: 1, Count: 0},
					{StartLine: 10, StartCol: 2, EndLine: 10, EndCol: 15, NumStmt: 1, Count: 0},
					{StartLine: 20, StartCol: 2, EndLine: 22, EndCol: 3, NumStmt: 2, Count: 0},
					{StartLine: 5, StartCol: 2, EndLine: 5, EndCol: 20, NumStmt: 1, Count: 1},
				},
			},
			expectCoverage:   0.2,
			expectUncovered:  3,
			expectLineRanges: []string{"10", "20-22", "30"},
		},
		{
			name: "coverage rate rounding to 4 decimals",
			profile: &cover.Profile{
				FileName: "rounding.go",
				Mode:     "set",
				Blocks: []cover.ProfileBlock{
					{StartLine: 5, StartCol: 2, EndLine: 5, EndCol: 20, NumStmt: 1, Count: 1},
					{StartLine: 10, StartCol: 2, EndLine: 10, EndCol: 15, NumStmt: 1, Count: 1},
					{StartLine: 15, StartCol: 2, EndLine: 15, EndCol: 10, NumStmt: 1, Count: 0},
				},
			},
			expectCoverage:   0.6667,
			expectUncovered:  1,
			expectLineRanges: []string{"15"},
		},
		{
			name: "complex coverage scenario",
			profile: &cover.Profile{
				FileName: "complex.go",
				Mode:     "set",
				Blocks: []cover.ProfileBlock{
					{StartLine: 1, StartCol: 1, EndLine: 1, EndCol: 10, NumStmt: 1, Count: 1},
					{StartLine: 5, StartCol: 2, EndLine: 8, EndCol: 3, NumStmt: 4, Count: 0},
					{StartLine: 10, StartCol: 2, EndLine: 10, EndCol: 15, NumStmt: 1, Count: 1},
					{StartLine: 12, StartCol: 2, EndLine: 12, EndCol: 20, NumStmt: 1, Count: 0},
					{StartLine: 15, StartCol: 2, EndLine: 18, EndCol: 3, NumStmt: 3, Count: 1},
					{StartLine: 20, StartCol: 2, EndLine: 25, EndCol: 3, NumStmt: 5, Count: 0},
				},
			},
			expectCoverage:   0.3333,
			expectUncovered:  3,
			expectLineRanges: []string{"5-8", "12", "20-25"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateFileCoverage(tt.profile)

			// Check file name
			if result.File != tt.profile.FileName {
				t.Errorf("File = %v, expect %v", result.File, tt.profile.FileName)
			}

			// Check coverage rate
			if result.CoverageRate != tt.expectCoverage {
				t.Errorf("CoverageRate = %v, expect %v", result.CoverageRate, tt.expectCoverage)
			}

			// Check number of uncovered items
			if len(result.UncoveredItems) != tt.expectUncovered {
				t.Errorf("UncoveredItems count = %v, expect %v", len(result.UncoveredItems), tt.expectUncovered)
			}

			// Check line ranges
			if len(result.UncoveredItems) > 0 {
				for i, item := range result.UncoveredItems {
					if i >= len(tt.expectLineRanges) {
						t.Errorf("Unexpected uncovered item at index %d: %v", i, item.LineRange)
						continue
					}
					if item.LineRange != tt.expectLineRanges[i] {
						t.Errorf("UncoveredItems[%d].LineRange = %v, expect %v", i, item.LineRange, tt.expectLineRanges[i])
					}
				}
			}

			// Check that enrichment fields are empty strings in Phase 1
			for i, item := range result.UncoveredItems {
				if item.Function != "" {
					t.Errorf("UncoveredItems[%d].Function = %q, expect empty string", i, item.Function)
				}
				if item.Type != "" {
					t.Errorf("UncoveredItems[%d].Type = %q, expect empty string", i, item.Type)
				}
				if item.Condition != "" {
					t.Errorf("UncoveredItems[%d].Condition = %q, expect empty string", i, item.Condition)
				}
			}
		})
	}
}

func TestFormatLineRange(t *testing.T) {
	tests := []struct {
		name      string
		startLine int
		endLine   int
		expect    string
	}{
		{
			name:      "single line",
			startLine: 10,
			endLine:   10,
			expect:    "10",
		},
		{
			name:      "multi-line range",
			startLine: 10,
			endLine:   15,
			expect:    "10-15",
		},
		{
			name:      "two line range",
			startLine: 5,
			endLine:   6,
			expect:    "5-6",
		},
		{
			name:      "large range",
			startLine: 100,
			endLine:   200,
			expect:    "100-200",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatLineRange(tt.startLine, tt.endLine)
			if got != tt.expect {
				t.Errorf("formatLineRange(%d, %d) = %v, expect %v", tt.startLine, tt.endLine, got, tt.expect)
			}
		})
	}
}

func TestExtractStartLine(t *testing.T) {
	tests := []struct {
		name      string
		lineRange string
		expect    int
	}{
		{
			name:      "single line",
			lineRange: "10",
			expect:    10,
		},
		{
			name:      "multi-line range",
			lineRange: "10-15",
			expect:    10,
		},
		{
			name:      "large line number",
			lineRange: "1000",
			expect:    1000,
		},
		{
			name:      "large range",
			lineRange: "500-1000",
			expect:    500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractStartLine(tt.lineRange)
			if got != tt.expect {
				t.Errorf("extractStartLine(%q) = %v, expect %v", tt.lineRange, got, tt.expect)
			}
		})
	}
}
