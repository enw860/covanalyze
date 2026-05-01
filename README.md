# covanalyze

A deterministic, offline Go CLI tool that parses coverage.out files and generates structured JSON reports with file-level coverage rates and uncovered line ranges.

## Overview

covanalyze is a Phase 1 implementation that reads Go coverage files (set mode only) using `golang.org/x/tools/cover`, calculates file-level coverage metrics, identifies uncovered code regions, and outputs structured JSON conforming to the Phase 0 schema.

### Features

- Parses standard Go coverage.out files (set mode only)
- Calculates file-level coverage rates with 4 decimal precision
- Identifies uncovered line ranges (one item per coverage block)
- Outputs pretty-printed JSON with 2-space indentation
- Configurable logging with glog
- Proper error handling with specific exit codes
- Works completely offline with no external dependencies at runtime

### Limitations (Phase 1)

- Only supports 'set' mode coverage files (count/atomic modes return error)
- No AST parsing or semantic enrichment
- No function name, type, or condition detection
- No range merging for adjacent/overlapping uncovered blocks

## Installation

### Prerequisites

- Go 1.24.3 or compatible version
- make (optional, for using Makefile targets)

### Build from Source

```bash
# Clone the repository
git clone <repository-url>
cd go_cov_analysis

# Build using make
make build

# Or build directly with go
go build -o bin/covanalyze ./cmd
```

The binary will be created at `bin/covanalyze`.

### Install

```bash
make install
```

This builds the binary to `bin/covanalyze` (no system-wide installation).

## Usage

### Basic Syntax

```bash
covanalyze [flags]
```

### Flags

- `-f <file>` - Coverage file path (required)
- `-o <file>` - Output file path (default: stdout)
- `-m <string>` - Module prefix to replace in file paths (optional, requires -mpath)
- `-mpath <string>` - Path prefix to replace module with (optional, requires -m)
- `-v <level>` - Verbose logging level (0=off, 1=info, 2=debug, 3=trace)
- `--help` - Show help message

### Examples

#### 1. Basic Usage

Parse a coverage file and output JSON to stdout:

```bash
covanalyze -f coverage.out
```

#### 2. Output to File

Write the JSON report to a file:

```bash
covanalyze -f coverage.out -o report.json
```

#### 3. Path Normalization

Normalize file paths by replacing module prefix with a local path:

```bash
covanalyze -f coverage.out -m github.com/user/repo -mpath ./
```

This is useful when you want to convert module-prefixed paths (e.g., `github.com/user/repo/pkg/file.go`) to local paths (e.g., `./pkg/file.go`). Both `-m` and `-mpath` flags must be provided together.

#### 4. Verbose Logging

Enable verbose logging to see detailed processing information:

```bash
covanalyze -f coverage.out -v 2
```

#### 5. Help Flag

Display usage information:

```bash
covanalyze --help
```

## Output Format

The tool outputs JSON conforming to the Phase 0 schema:

```json
{
  "$schema": "./schema/output.schema.json",
  "file_reports": [
    {
      "file": "path/to/file.go",
      "coverage_rate": 0.7500,
      "uncovered_items": [
        {
          "function": "",
          "type": "",
          "condition": "",
          "line_range": "10-15"
        }
      ]
    }
  ]
}
```

### Output Details

- **$schema**: Reference to the JSON schema (hardcoded to `./schema/output.schema.json`)
- **file_reports**: Array of file coverage reports, sorted by coverage_rate descending
- **file**: File path exactly as it appears in coverage.out (no normalization)
- **coverage_rate**: Coverage percentage as decimal (0.0 to 1.0), rounded to 4 decimal places
- **uncovered_items**: Array of uncovered code segments, sorted by line number ascending
  - **function**: Empty string in Phase 1 (enrichment in Phase 2)
  - **type**: Empty string in Phase 1 (enrichment in Phase 2)
  - **condition**: Empty string in Phase 1 (enrichment in Phase 2)
  - **line_range**: Line range in format "N" (single line) or "N-M" (multi-line)

### Special Cases

- **Empty coverage file**: Returns empty `file_reports` array with exit code 0
- **Files with 100% coverage**: Included with empty `uncovered_items` array
- **Files with 0 statements**: Coverage rate is 1.0

## Exit Codes

The tool uses specific exit codes to indicate different error conditions:

| Exit Code | Description |
|-----------|-------------|
| 0 | Success - coverage file parsed and JSON generated successfully |
| 1 | File not found - the specified coverage file does not exist |
| 2 | Parse error - invalid coverage file format or unsupported mode (count/atomic) |
| 3 | Validation error - internal validation or JSON generation failed |

### Error Handling

All errors are written to stderr, while JSON output goes to stdout (or the file specified with `-o`).

Example error messages:

```bash
# File not found
$ covanalyze -f nonexistent.out
Error: file not found: nonexistent.out

# Unsupported mode
$ covanalyze -f count_mode.out
Error: unsupported coverage mode 'count': only 'set' mode is supported

# Missing required flag
$ covanalyze
Error: coverage file path is required (use -f flag)
Run 'covanalyze --help' for usage information
```

## Development

### Running Tests

```bash
# Run all tests with coverage
make test

# Or run tests directly
go test -v -cover ./...
```

### Building

```bash
# Build the binary
make build

# Clean build artifacts
make clean
```

### Project Structure

```
.
├── cmd/
│   ├── main.go           # CLI entry point
│   └── usage.txt         # Help text
├── internal/
│   ├── errors/           # Custom error types
│   ├── formatter/        # JSON formatting
│   ├── models/           # Data structures
│   └── parser/           # Coverage parsing and calculation
├── schema/
│   └── output.schema.json # JSON schema
├── test/
│   └── resources/        # Test fixtures
├── Makefile              # Build automation
└── README.md             # This file
```

## Schema Reference

The output JSON conforms to the schema defined in `schema/output.schema.json`. This schema is part of Phase 0 and defines the structure for coverage analysis reports.

## License

See LICENSE file for details.

## Contributing

This is Phase 1 of the covanalyze project. Future phases will add:
- Phase 2: AST parsing and semantic enrichment (function names, types, conditions)
- Phase 3: Additional coverage modes and advanced features

For questions or issues, please refer to the project documentation.