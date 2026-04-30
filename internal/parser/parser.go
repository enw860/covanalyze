package parser

import (
	"os"

	"golang.org/x/tools/cover"

	"github.com/enw860/covanalyze/internal/errors"
)

const (
	// SupportedMode is the only coverage mode supported in Phase 1
	SupportedMode = "set"
)

// ParseCoverageFile parses a coverage.out file and returns the coverage profiles.
// It validates that the file exists and uses 'set' mode.
// Returns FileNotFoundError if the file doesn't exist.
// Returns UnsupportedModeError if the coverage mode is not 'set'.
func ParseCoverageFile(path string) ([]*cover.Profile, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, &errors.FileNotFoundError{Path: path}
	}

	// Parse the coverage file
	profiles, err := cover.ParseProfiles(path)
	if err != nil {
		// Wrap parse errors as UnsupportedModeError since golang.org/x/tools/cover
		// returns generic errors for parsing issues
		return nil, &errors.UnsupportedModeError{
			Mode:    "unknown",
			Message: err.Error(),
		}
	}

	// Validate that all profiles use 'set' mode
	for _, profile := range profiles {
		if profile.Mode != SupportedMode {
			return nil, &errors.UnsupportedModeError{
				Mode:    profile.Mode,
				Message: "only 'set' mode is supported",
			}
		}
	}

	return profiles, nil
}

// Made with Bob
