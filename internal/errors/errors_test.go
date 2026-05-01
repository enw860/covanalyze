package errors

import (
	"fmt"
	"testing"
)

func TestFileNotFoundError(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{
			name: "coverage file",
			path: "coverage.out",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &FileNotFoundError{Path: tt.path}
			expected := fmt.Sprintf(fileNotFoundTemplate, tt.path)
			if got := err.Error(); got != expected {
				t.Errorf("FileNotFoundError.Error() = %q, want %q", got, expected)
			}
		})
	}
}

func TestUnsupportedModeError(t *testing.T) {
	tests := []struct {
		name    string
		mode    string
		message string
	}{
		{
			name:    "count mode without message",
			mode:    "count",
			message: "",
		},
		{
			name:    "atomic mode without message",
			mode:    "atomic",
			message: "",
		},
		{
			name:    "count mode with custom message",
			mode:    "count",
			message: "only supports set mode",
		},
		{
			name:    "atomic mode with custom message",
			mode:    "atomic",
			message: "use set mode instead",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &UnsupportedModeError{Mode: tt.mode, Message: tt.message}
			msg := tt.message
			if msg == "" {
				msg = unsupportedModeDefaultMsg
			}
			expected := fmt.Sprintf(unsupportedModeTemplate, tt.mode, msg)
			if got := err.Error(); got != expected {
				t.Errorf("UnsupportedModeError.Error() = %q, want %q", got, expected)
			}
		})
	}
}
