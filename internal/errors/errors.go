package errors

import "fmt"

const (
	// Error message templates
	fileNotFoundTemplate      = "file not found: %s"
	unsupportedModeTemplate   = "unsupported coverage mode '%s': %s"
	unsupportedModeDefaultMsg = "only 'set' mode is supported"
)

// FileNotFoundError represents an error when a file cannot be found.
// This error should result in exit code 1.
type FileNotFoundError struct {
	Path string
}

// Error implements the error interface for FileNotFoundError.
func (e *FileNotFoundError) Error() string {
	return fmt.Sprintf(fileNotFoundTemplate, e.Path)
}

// UnsupportedModeError represents an error when a coverage file uses an unsupported mode.
// Only "set" mode is supported in Phase 1. This error should result in exit code 2.
type UnsupportedModeError struct {
	Mode    string
	Message string
}

// Error implements the error interface for UnsupportedModeError.
func (e *UnsupportedModeError) Error() string {
	msg := e.Message
	if msg == "" {
		msg = unsupportedModeDefaultMsg
	}
	return fmt.Sprintf(unsupportedModeTemplate, e.Mode, msg)
}
