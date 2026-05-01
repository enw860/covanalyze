package parser

import (
	"strings"
	"testing"

	"github.com/enw860/covanalyze/internal/models"
)

func TestParseSourceFile_IfStatements(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		expected []models.ASTNode
	}{
		{
			name: "simple if statement",
			file: "../../test/enricher/if_simple.go",
			expected: []models.ASTNode{
				{
					NodeType:     models.NodeTypeIf,
					StartLine:    4,
					EndLine:      6,
					FunctionName: "ifSimple",
					Condition:    "x > 0",
				},
			},
		},
		{
			name: "if-else statement",
			file: "../../test/enricher/if_else.go",
			expected: []models.ASTNode{
				{
					NodeType:     models.NodeTypeElse,
					StartLine:    6,
					EndLine:      8,
					FunctionName: "ifElse",
					Condition:    "",
				},
				{
					NodeType:     models.NodeTypeIf,
					StartLine:    4,
					EndLine:      8,
					FunctionName: "ifElse",
					Condition:    "x > 0",
				},
			},
		},
		{
			name: "if-else-if statement",
			file: "../../test/enricher/if_elseif.go",
			expected: []models.ASTNode{
				{
					NodeType:     models.NodeTypeIf,
					StartLine:    4,
					EndLine:      8,
					FunctionName: "ifElseIf",
					Condition:    "x > 0",
				},
				{
					NodeType:     models.NodeTypeIf,
					StartLine:    6,
					EndLine:      8,
					FunctionName: "ifElseIf",
					Condition:    "x < 0",
				},
			},
		},
		{
			name: "nested if statements",
			file: "../../test/enricher/if_nested.go",
			expected: []models.ASTNode{
				{
					NodeType:     models.NodeTypeIf,
					StartLine:    4,
					EndLine:      8,
					FunctionName: "ifNested",
					Condition:    "x > 0",
				},
				{
					NodeType:     models.NodeTypeIf,
					StartLine:    5,
					EndLine:      7,
					FunctionName: "ifNested",
					Condition:    "y > 0",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the file
			nodes, err := parseSourceFile(tt.file)
			if err != nil {
				t.Fatalf("parseSourceFile() error = %v", err)
			}

			// Verify node count
			if len(nodes) != len(tt.expected) {
				t.Errorf("parseSourceFile() returned %d nodes, expected %d", len(nodes), len(tt.expected))
			}

			// Verify each node
			for i, expected := range tt.expected {
				if i >= len(nodes) {
					break
				}
				node := nodes[i]
				if node.NodeType != expected.NodeType {
					t.Errorf("Node[%d].NodeType = %q, expected %q", i, node.NodeType, expected.NodeType)
				}
				if node.StartLine != expected.StartLine {
					t.Errorf("Node[%d].StartLine = %d, expected %d", i, node.StartLine, expected.StartLine)
				}
				if node.EndLine != expected.EndLine {
					t.Errorf("Node[%d].EndLine = %d, expected %d", i, node.EndLine, expected.EndLine)
				}
				if node.FunctionName != expected.FunctionName {
					t.Errorf("Node[%d].FunctionName = %q, expected %q", i, node.FunctionName, expected.FunctionName)
				}
				if node.Condition != expected.Condition {
					t.Errorf("Node[%d].Condition = %q, expected %q", i, node.Condition, expected.Condition)
				}
			}
		})
	}
}

func TestParseSourceFile_SwitchStatements(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		expected []models.ASTNode
	}{
		{
			name: "simple switch statement",
			file: "../../test/enricher/switch_simple.go",
			expected: []models.ASTNode{
				{
					NodeType:     models.NodeTypeSwitch,
					StartLine:    4,
					EndLine:      9,
					FunctionName: "switchSimple",
					Condition:    "x",
				},
				{
					NodeType:     models.NodeTypeCase,
					StartLine:    5,
					EndLine:      6,
					FunctionName: "switchSimple",
					Condition:    "1",
				},
				{
					NodeType:     models.NodeTypeCase,
					StartLine:    7,
					EndLine:      8,
					FunctionName: "switchSimple",
					Condition:    "2",
				},
			},
		},
		{
			name: "switch with multiple case values",
			file: "../../test/enricher/switch_multiple.go",
			expected: []models.ASTNode{
				{
					NodeType:     models.NodeTypeSwitch,
					StartLine:    4,
					EndLine:      7,
					FunctionName: "switchMultiple",
					Condition:    "x",
				},
				{
					NodeType:     models.NodeTypeCase,
					StartLine:    5,
					EndLine:      6,
					FunctionName: "switchMultiple",
					Condition:    "1, 2, 3",
				},
			},
		},
		{
			name: "switch with default case",
			file: "../../test/enricher/switch_default.go",
			expected: []models.ASTNode{
				{
					NodeType:     models.NodeTypeSwitch,
					StartLine:    4,
					EndLine:      9,
					FunctionName: "switchDefault",
					Condition:    "x",
				},
				{
					NodeType:     models.NodeTypeCase,
					StartLine:    5,
					EndLine:      6,
					FunctionName: "switchDefault",
					Condition:    "1",
				},
				{
					NodeType:     models.NodeTypeCase,
					StartLine:    7,
					EndLine:      8,
					FunctionName: "switchDefault",
					Condition:    "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, err := parseSourceFile(tt.file)
			if err != nil {
				t.Fatalf("parseSourceFile() error = %v", err)
			}

			if len(nodes) != len(tt.expected) {
				t.Errorf("parseSourceFile() returned %d nodes, expected %d", len(nodes), len(tt.expected))
			}

			for i, expected := range tt.expected {
				if i >= len(nodes) {
					break
				}
				node := nodes[i]
				if node.NodeType != expected.NodeType {
					t.Errorf("Node[%d].NodeType = %q, expected %q", i, node.NodeType, expected.NodeType)
				}
				if node.StartLine != expected.StartLine {
					t.Errorf("Node[%d].StartLine = %d, expected %d", i, node.StartLine, expected.StartLine)
				}
				if node.EndLine != expected.EndLine {
					t.Errorf("Node[%d].EndLine = %d, expected %d", i, node.EndLine, expected.EndLine)
				}
				if node.FunctionName != expected.FunctionName {
					t.Errorf("Node[%d].FunctionName = %q, expected %q", i, node.FunctionName, expected.FunctionName)
				}
				if node.Condition != expected.Condition {
					t.Errorf("Node[%d].Condition = %q, expected %q", i, node.Condition, expected.Condition)
				}
			}
		})
	}
}

func TestParseSourceFile_TypeSwitchStatements(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		expected []models.ASTNode
	}{
		{
			name: "type switch statement",
			file: "../../test/enricher/typeswitch.go",
			expected: []models.ASTNode{
				{
					NodeType:     models.NodeTypeTypeSwitch,
					StartLine:    4,
					EndLine:      12,
					FunctionName: "typeSwitch",
					Condition:    "",
				},
				{
					NodeType:     models.NodeTypeCase,
					StartLine:    5,
					EndLine:      6,
					FunctionName: "typeSwitch",
					Condition:    "int",
				},
				{
					NodeType:     models.NodeTypeCase,
					StartLine:    7,
					EndLine:      8,
					FunctionName: "typeSwitch",
					Condition:    "string",
				},
				{
					NodeType:     models.NodeTypeCase,
					StartLine:    9,
					EndLine:      11,
					FunctionName: "typeSwitch",
					Condition:    "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, err := parseSourceFile(tt.file)
			if err != nil {
				t.Fatalf("parseSourceFile() error = %v", err)
			}

			if len(nodes) != len(tt.expected) {
				t.Errorf("parseSourceFile() returned %d nodes, expected %d", len(nodes), len(tt.expected))
			}

			for i, expected := range tt.expected {
				if i >= len(nodes) {
					break
				}
				node := nodes[i]
				if node.NodeType != expected.NodeType {
					t.Errorf("Node[%d].NodeType = %q, expected %q", i, node.NodeType, expected.NodeType)
				}
				if node.StartLine != expected.StartLine {
					t.Errorf("Node[%d].StartLine = %d, expected %d", i, node.StartLine, expected.StartLine)
				}
				if node.EndLine != expected.EndLine {
					t.Errorf("Node[%d].EndLine = %d, expected %d", i, node.EndLine, expected.EndLine)
				}
				if node.FunctionName != expected.FunctionName {
					t.Errorf("Node[%d].FunctionName = %q, expected %q", i, node.FunctionName, expected.FunctionName)
				}
			}
		})
	}
}

func TestParseSourceFile_ForLoops(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		expected []models.ASTNode
	}{
		{
			name: "for loop with condition",
			file: "../../test/enricher/for_condition.go",
			expected: []models.ASTNode{
				{
					NodeType:     models.NodeTypeFor,
					StartLine:    4,
					EndLine:      6,
					FunctionName: "forCondition",
					Condition:    "i < 10",
				},
			},
		},
		{
			name: "for loop with init, condition, post",
			file: "../../test/enricher/for_full.go",
			expected: []models.ASTNode{
				{
					NodeType:     models.NodeTypeFor,
					StartLine:    5,
					EndLine:      7,
					FunctionName: "forFull",
					Condition:    "i < 10",
				},
			},
		},
		{
			name: "infinite for loop",
			file: "../../test/enricher/for_infinite.go",
			expected: []models.ASTNode{
				{
					NodeType:     models.NodeTypeFor,
					StartLine:    5,
					EndLine:      10,
					FunctionName: "forInfinite",
					Condition:    "",
				},
				{
					NodeType:     models.NodeTypeIf,
					StartLine:    7,
					EndLine:      9,
					FunctionName: "forInfinite",
					Condition:    "count > 5",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, err := parseSourceFile(tt.file)
			if err != nil {
				t.Fatalf("parseSourceFile() error = %v", err)
			}

			if len(nodes) != len(tt.expected) {
				t.Errorf("parseSourceFile() returned %d nodes, expected %d", len(nodes), len(tt.expected))
			}

			for i, expected := range tt.expected {
				if i >= len(nodes) {
					break
				}
				node := nodes[i]
				if node.NodeType != expected.NodeType {
					t.Errorf("Node[%d].NodeType = %q, expected %q", i, node.NodeType, expected.NodeType)
				}
				if node.StartLine != expected.StartLine {
					t.Errorf("Node[%d].StartLine = %d, expected %d", i, node.StartLine, expected.StartLine)
				}
				if node.EndLine != expected.EndLine {
					t.Errorf("Node[%d].EndLine = %d, expected %d", i, node.EndLine, expected.EndLine)
				}
				if node.FunctionName != expected.FunctionName {
					t.Errorf("Node[%d].FunctionName = %q, expected %q", i, node.FunctionName, expected.FunctionName)
				}
				if node.Condition != expected.Condition {
					t.Errorf("Node[%d].Condition = %q, expected %q", i, node.Condition, expected.Condition)
				}
			}
		})
	}
}

func TestParseSourceFile_RangeLoops(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		expected []models.ASTNode
	}{
		{
			name: "range over slice",
			file: "../../test/enricher/range_slice.go",
			expected: []models.ASTNode{
				{
					NodeType:     models.NodeTypeRange,
					StartLine:    5,
					EndLine:      7,
					FunctionName: "rangeSlice",
					Condition:    "items",
				},
			},
		},
		{
			name: "range with blank identifier",
			file: "../../test/enricher/range_blank.go",
			expected: []models.ASTNode{
				{
					NodeType:     models.NodeTypeRange,
					StartLine:    5,
					EndLine:      7,
					FunctionName: "rangeBlank",
					Condition:    "items",
				},
			},
		},
		{
			name: "range over map",
			file: "../../test/enricher/range_map.go",
			expected: []models.ASTNode{
				{
					NodeType:     models.NodeTypeRange,
					StartLine:    5,
					EndLine:      7,
					FunctionName: "rangeMap",
					Condition:    "myMap",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, err := parseSourceFile(tt.file)
			if err != nil {
				t.Fatalf("parseSourceFile() error = %v", err)
			}

			if len(nodes) != len(tt.expected) {
				t.Errorf("parseSourceFile() returned %d nodes, expected %d", len(nodes), len(tt.expected))
			}

			for i, expected := range tt.expected {
				if i >= len(nodes) {
					break
				}
				node := nodes[i]
				if node.NodeType != expected.NodeType {
					t.Errorf("Node[%d].NodeType = %q, expected %q", i, node.NodeType, expected.NodeType)
				}
				if node.StartLine != expected.StartLine {
					t.Errorf("Node[%d].StartLine = %d, expected %d", i, node.StartLine, expected.StartLine)
				}
				if node.EndLine != expected.EndLine {
					t.Errorf("Node[%d].EndLine = %d, expected %d", i, node.EndLine, expected.EndLine)
				}
				if node.FunctionName != expected.FunctionName {
					t.Errorf("Node[%d].FunctionName = %q, expected %q", i, node.FunctionName, expected.FunctionName)
				}
				if node.Condition != expected.Condition {
					t.Errorf("Node[%d].Condition = %q, expected %q", i, node.Condition, expected.Condition)
				}
			}
		})
	}
}

func TestParseSourceFile_MissingFile(t *testing.T) {
	nodes, err := parseSourceFile("/nonexistent/file.go")
	if err == nil {
		t.Error("parseSourceFile() expected error for missing file, got nil")
	}
	if nodes != nil {
		t.Errorf("parseSourceFile() expected nil nodes for missing file, got %v", nodes)
	}
}

func TestEnrichFileReports(t *testing.T) {
	longFunctionName := strings.Repeat("f", models.MaxFunctionNameLength+10)
	longCondition := strings.Repeat("c", models.MaxConditionLength+10)

	tests := []struct {
		name     string
		reports  []models.FileReport
		validate func(t *testing.T, reports []models.FileReport)
	}{
		{
			name: "enriches items in place",
			reports: []models.FileReport{
				{
					File: "../../test/enricher/if_simple.go",
					UncoveredItems: []models.UncoveredItem{
						{LineRange: "4-6"},
					},
				},
				{
					File: "../../test/enricher/for_condition.go",
					UncoveredItems: []models.UncoveredItem{
						{LineRange: "4-6"},
					},
				},
				{
					File: "../../test/enricher/switch_simple.go",
					UncoveredItems: []models.UncoveredItem{
						{LineRange: "5"},
					},
				},
			},
			validate: func(t *testing.T, reports []models.FileReport) {
				t.Helper()

				if reports[0].UncoveredItems[0].Function != "ifSimple" {
					t.Errorf("Function = %q, expected %q", reports[0].UncoveredItems[0].Function, "ifSimple")
				}
				if reports[0].UncoveredItems[0].Type != "branch" {
					t.Errorf("Type = %q, expected %q", reports[0].UncoveredItems[0].Type, "branch")
				}
				if reports[0].UncoveredItems[0].Condition != "x > 0" {
					t.Errorf("Condition = %q, expected %q", reports[0].UncoveredItems[0].Condition, "x > 0")
				}
				if reports[0].UncoveredItems[0].LineRange != "4-6" {
					t.Errorf("LineRange = %q, expected %q", reports[0].UncoveredItems[0].LineRange, "4-6")
				}

				if reports[1].UncoveredItems[0].Function != "forCondition" {
					t.Errorf("Function = %q, expected %q", reports[1].UncoveredItems[0].Function, "forCondition")
				}
				if reports[1].UncoveredItems[0].Type != "loop" {
					t.Errorf("Type = %q, expected %q", reports[1].UncoveredItems[0].Type, "loop")
				}
				if reports[1].UncoveredItems[0].Condition != "i < 10" {
					t.Errorf("Condition = %q, expected %q", reports[1].UncoveredItems[0].Condition, "i < 10")
				}
				if reports[1].UncoveredItems[0].LineRange != "4-6" {
					t.Errorf("LineRange = %q, expected %q", reports[1].UncoveredItems[0].LineRange, "4-6")
				}

				if reports[2].UncoveredItems[0].Function != "switchSimple" {
					t.Errorf("Function = %q, expected %q", reports[2].UncoveredItems[0].Function, "switchSimple")
				}
				if reports[2].UncoveredItems[0].Type != "branch" {
					t.Errorf("Type = %q, expected %q", reports[2].UncoveredItems[0].Type, "branch")
				}
				if reports[2].UncoveredItems[0].Condition != "1" {
					t.Errorf("Condition = %q, expected %q", reports[2].UncoveredItems[0].Condition, "1")
				}
				if reports[2].UncoveredItems[0].LineRange != "5" {
					t.Errorf("LineRange = %q, expected %q", reports[2].UncoveredItems[0].LineRange, "5")
				}
			},
		},
		{
			name: "leaves fields empty for missing file",
			reports: []models.FileReport{
				{
					File: "/nonexistent/file.go",
					UncoveredItems: []models.UncoveredItem{
						{LineRange: "10-12"},
					},
				},
			},
			validate: func(t *testing.T, reports []models.FileReport) {
				t.Helper()

				item := reports[0].UncoveredItems[0]
				if item.Function != "" {
					t.Errorf("Function = %q, expected empty string", item.Function)
				}
				if item.Type != "" {
					t.Errorf("Type = %q, expected empty string", item.Type)
				}
				if item.Condition != "" {
					t.Errorf("Condition = %q, expected empty string", item.Condition)
				}
				if item.LineRange != "10-12" {
					t.Errorf("LineRange = %q, expected %q", item.LineRange, "10-12")
				}
			},
		},
		{
			name: "truncates function and condition values",
			reports: []models.FileReport{
				{
					File: "../../test/enricher/if_simple.go",
					UncoveredItems: []models.UncoveredItem{
						{LineRange: "4-6"},
					},
				},
			},
			validate: func(t *testing.T, reports []models.FileReport) {
				t.Helper()

				reports[0].UncoveredItems[0].Function = truncateString(longFunctionName, models.MaxFunctionNameLength)
				reports[0].UncoveredItems[0].Condition = truncateString(longCondition, models.MaxConditionLength)

				if len(reports[0].UncoveredItems[0].Function) != models.MaxFunctionNameLength {
					t.Errorf("Function length = %d, expected %d", len(reports[0].UncoveredItems[0].Function), models.MaxFunctionNameLength)
				}
				if !strings.HasSuffix(reports[0].UncoveredItems[0].Function, "...") {
					t.Errorf("Function = %q, expected ellipsis suffix", reports[0].UncoveredItems[0].Function)
				}
				if len(reports[0].UncoveredItems[0].Condition) != models.MaxConditionLength {
					t.Errorf("Condition length = %d, expected %d", len(reports[0].UncoveredItems[0].Condition), models.MaxConditionLength)
				}
				if !strings.HasSuffix(reports[0].UncoveredItems[0].Condition, "...") {
					t.Errorf("Condition = %q, expected ellipsis suffix", reports[0].UncoveredItems[0].Condition)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			EnrichFileReports(tt.reports)
			tt.validate(t, tt.reports)
		})
	}
}
