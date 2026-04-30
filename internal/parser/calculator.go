package parser

import (
	"fmt"
	"math"
	"sort"

	"golang.org/x/tools/cover"

	"github.com/enw860/covanalyze/internal/models"
)

// CalculateFileCoverage calculates coverage metrics for a single file profile.
// It computes the coverage rate, identifies uncovered blocks, and returns a FileReport.
// Coverage rate is rounded to 4 decimal places.
// Files with 0 statements return coverage_rate = 1.0.
// Each uncovered block becomes a separate UncoveredItem (no merging).
// UncoveredItems are sorted by line number ascending.
// Function, type, and condition fields are empty strings in Phase 1.
func CalculateFileCoverage(profile *cover.Profile) models.FileReport {
	var totalStatements int
	var coveredStatements int
	var uncoveredItems []models.UncoveredItem

	// Process each block in the profile
	for _, block := range profile.Blocks {
		totalStatements += block.NumStmt

		if block.Count > 0 {
			// Block is covered
			coveredStatements += block.NumStmt
		} else {
			// Block is uncovered - create an UncoveredItem
			lineRange := formatLineRange(block.StartLine, block.EndLine)
			uncoveredItems = append(uncoveredItems, models.UncoveredItem{
				Function:  "", // Empty in Phase 1
				Type:      "", // Empty in Phase 1
				Condition: "", // Empty in Phase 1
				LineRange: lineRange,
			})
		}
	}

	// Calculate coverage rate
	var coverageRate float64
	if totalStatements == 0 {
		// Files with 0 statements have coverage_rate = 1.0
		coverageRate = 1.0
	} else {
		coverageRate = float64(coveredStatements) / float64(totalStatements)
		// Round to 4 decimal places
		coverageRate = math.Round(coverageRate*10000) / 10000
	}

	// Sort uncovered items by line number ascending
	sort.Slice(uncoveredItems, func(i, j int) bool {
		return extractStartLine(uncoveredItems[i].LineRange) < extractStartLine(uncoveredItems[j].LineRange)
	})

	return models.FileReport{
		File:           profile.FileName,
		CoverageRate:   coverageRate,
		UncoveredItems: uncoveredItems,
	}
}

// formatLineRange formats a line range as 'N' for single line or 'N-M' for multi-line.
func formatLineRange(startLine, endLine int) string {
	if startLine == endLine {
		return fmt.Sprintf("%d", startLine)
	}
	return fmt.Sprintf("%d-%d", startLine, endLine)
}

// extractStartLine extracts the start line number from a line range string.
// Used for sorting UncoveredItems.
func extractStartLine(lineRange string) int {
	var startLine int
	fmt.Sscanf(lineRange, "%d", &startLine)
	return startLine
}
