package parser

import (
	"path/filepath"
	"testing"

	"github.com/enw860/covanalyze/internal/errors"
)

func TestParseCoverageFile(t *testing.T) {
	tests := []struct {
		name           string
		fixturePath    string
		expectErr      bool
		expectErrType  string
		expectProfiles int
		checkMode      bool
	}{
		{
			name:           "successful set mode parsing",
			fixturePath:    "test/resources/coverage/valid_set_mode.out",
			expectErr:      false,
			expectProfiles: 1,
			checkMode:      true,
		},
		{
			name:          "file not found error",
			fixturePath:   "test/resources/coverage/nonexistent.out",
			expectErr:     true,
			expectErrType: "*errors.FileNotFoundError",
		},
		{
			name:          "invalid coverage format",
			fixturePath:   "test/resources/coverage/invalid_format.out",
			expectErr:     true,
			expectErrType: "*errors.UnsupportedModeError",
		},
		{
			name:          "count mode rejection",
			fixturePath:   "test/resources/coverage/count_mode.out",
			expectErr:     true,
			expectErrType: "*errors.UnsupportedModeError",
		},
		{
			name:          "atomic mode rejection",
			fixturePath:   "test/resources/coverage/atomic_mode.out",
			expectErr:     true,
			expectErrType: "*errors.UnsupportedModeError",
		},
		{
			name:           "empty coverage file with mode line only",
			fixturePath:    "test/resources/coverage/empty_set_mode.out",
			expectErr:      false,
			expectProfiles: 0,
		},
		{
			name:           "multiple files in set mode",
			fixturePath:    "test/resources/coverage/multiple_files_set_mode.out",
			expectErr:      false,
			expectProfiles: 3,
			checkMode:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get absolute path to fixture
			path, err := filepath.Abs(filepath.Join("../..", tt.fixturePath))
			if err != nil {
				t.Fatalf("failed to get absolute path: %v", err)
			}

			profiles, err := ParseCoverageFile(path)

			if tt.expectErr {
				if err == nil {
					t.Errorf("ParseCoverageFile() expected error but got none")
					return
				}

				// Check error type
				switch tt.expectErrType {
				case "*errors.FileNotFoundError":
					if _, ok := err.(*errors.FileNotFoundError); !ok {
						t.Errorf("ParseCoverageFile() error type = %T, want *errors.FileNotFoundError", err)
					}
				case "*errors.UnsupportedModeError":
					if _, ok := err.(*errors.UnsupportedModeError); !ok {
						t.Errorf("ParseCoverageFile() error type = %T, want *errors.UnsupportedModeError", err)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("ParseCoverageFile() unexpected error = %v", err)
				return
			}

			if len(profiles) != tt.expectProfiles {
				t.Errorf("ParseCoverageFile() got %d profiles, want %d", len(profiles), tt.expectProfiles)
			}

			if tt.checkMode && len(profiles) > 0 {
				for i, profile := range profiles {
					if profile.Mode != SupportedMode {
						t.Errorf("ParseCoverageFile() profile[%d].Mode = %s, want %s", i, profile.Mode, SupportedMode)
					}
				}
			}
		})
	}
}
