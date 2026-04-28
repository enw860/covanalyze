# Covanalyze Output Schema Documentation

## Overview

The Covanalyze output schema defines the JSON structure for coverage analysis results. It supports two operational modes:

1. **Coverage-only mode**: Basic coverage data derived solely from `coverage.out` files
2. **Enriched mode**: Enhanced coverage data with AST (Abstract Syntax Tree) analysis

The schema is defined using JSON Schema Draft 2020-12 and is located at [`schema/output.schema.json`](../schema/output.schema.json).

## Schema Metadata

- **$schema**: `https://json-schema.org/draft/2020-12/schema`
- **$id**: `covanalyze-output-schema-v1`
- **Compliance**: JSON Schema Draft 2020-12

## Root Output Structure

The root output object contains:

```json
{
  "$schema": "./schema/output.schema.json",
  "file_reports": [...]
}
```

### Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `$schema` | string | Yes | Reference to the JSON Schema file. Configurable, defaults to `./schema/output.schema.json` |
| `file_reports` | array | Yes | Array of FileReport objects, one per analyzed file |

### Constraints

- **Sorting**: FileReports are sorted by `coverage_rate` in descending order (highest coverage first)
- **Inclusion**: All files are included, even those with 100% coverage
- **Format**: JSON output is pretty-printed with 2-space indentation by default

## Data Types

### FileReport

Represents coverage analysis results for a single source file.

#### Structure

```json
{
  "file": "internal/processor/handler.go",
  "coverage_rate": 0.7250,
  "uncovered_items": [...]
}
```

#### Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `file` | string | Yes | File path exactly as it appears in coverage.out, relative to project root |
| `coverage_rate` | number | Yes | Coverage rate between 0.0 and 1.0, rounded to 4 decimal places |
| `uncovered_items` | array | Yes | Array of UncoveredItem objects, sorted by line number ascending |

#### Constraints

- **Path handling**: File paths are output exactly as in coverage.out with no normalization
- **Coverage calculation**: 
  - Rounded to 4 decimal places (e.g., `0.7845`)
  - Set to `1.0` for files with 0 statements
  - Range: `0.0` to `1.0` inclusive
- **Empty arrays**: Files with 100% coverage have an empty `uncovered_items` array

#### Example

```json
{
  "file": "internal/utils/helpers.go",
  "coverage_rate": 0.8500,
  "uncovered_items": [
    {
      "line_range": "42"
    }
  ]
}
```

### UncoveredItem

Represents a single uncovered code segment with optional AST enrichment data.

#### Structure

```json
{
  "function": "ProcessRequest",
  "type": "branch",
  "condition": "if err != nil",
  "line_range": "34-36"
}
```

#### Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `line_range` | string | Yes | Line range in format `"start"` or `"start-end"` |
| `function` | string | No | Name of the function containing this uncovered code |
| `type` | string | No | Type of uncovered code segment (enum) |
| `condition` | string | No | Condition expression for branches or loops |

#### Field Details

**line_range**
- Format: Single line (`"45"`) or range (`"45-52"`)
- Pattern: `^[1-9][0-9]*(-[1-9][0-9]*)?$`
- Validation: start_line > 0, end_line >= start_line
- 1-based line numbers matching coverage.out format

**function** (optional, AST enrichment)
- Truncated at 256 characters with `...` suffix if longer
- Uses `omitempty` JSON tag (omitted when empty)

**type** (optional, AST enrichment)
- Enum values:
  - `"branch"`: Conditional branch
  - `"error_path"`: Error handling path
  - `"early_return"`: Early return statement
  - `"loop"`: Loop body
  - `""`: Empty string when type cannot be determined
- Uses `omitempty` JSON tag (omitted when empty)

**condition** (optional, AST enrichment)
- Truncated at 512 characters with `...` suffix if longer
- Uses empty string when unavailable (not `"unknown_condition"`)
- Uses `omitempty` JSON tag (omitted when empty)

#### Constraints

- **Sorting**: Within a FileReport, uncovered_items are sorted by line number ascending
- **Optional fields**: `function`, `type`, and `condition` are populated only during AST enrichment

#### Examples

**Coverage-only mode** (minimal):
```json
{
  "line_range": "23-27"
}
```

**Enriched mode** (with AST data):
```json
{
  "function": "ValidateInput",
  "type": "early_return",
  "condition": "if input == nil",
  "line_range": "58-59"
}
```

### CoverageBlock

Represents a coverage block from coverage.out format. Used internally for parsing.

#### Structure

```go
type CoverageBlock struct {
    File      string `json:"file"`
    StartLine int    `json:"start_line"`
    StartCol  int    `json:"start_col"`
    EndLine   int    `json:"end_line"`
    EndCol    int    `json:"end_col"`
    NumStmt   int    `json:"num_stmt"`
    Count     int    `json:"count"`
}
```

#### Fields

| Field | Type | Description |
|-------|------|-------------|
| `file` | string | File path from coverage.out |
| `start_line` | int | Starting line number (1-based) |
| `start_col` | int | Starting column number (1-based) |
| `end_line` | int | Ending line number (1-based) |
| `end_col` | int | Ending column number (1-based) |
| `num_stmt` | int | Number of statements in block |
| `count` | int | Execution count (0 = uncovered) |

#### Constraints

- All positions are 1-based, matching coverage.out format
- Used for internal processing, not included in final JSON output

### ASTNode

Represents minimal AST information for code enrichment. Used internally for AST analysis.

#### Structure

```go
type ASTNode struct {
    NodeType      string `json:"node_type"`
    StartLine     int    `json:"start_line"`
    EndLine       int    `json:"end_line"`
    FunctionName  string `json:"function_name,omitempty"`
    ConditionText string `json:"condition_text,omitempty"`
}
```

#### Fields

| Field | Type | Description |
|-------|------|-------------|
| `node_type` | string | Type of AST node (e.g., "IfStmt", "ForStmt") |
| `start_line` | int | Starting line number |
| `end_line` | int | Ending line number |
| `function_name` | string | Name of containing function (optional) |
| `condition_text` | string | Condition expression text (optional) |

#### Constraints

- Minimal fields only - no column positions or parent references
- Optional fields use `omitempty` JSON tag
- Used for internal processing, not included in final JSON output

## Operational Modes

### Coverage-Only Mode

Uses only data from `coverage.out` files. No AST parsing required.

**Characteristics:**
- Only `line_range` field populated in UncoveredItem
- No `function`, `type`, or `condition` fields
- Faster processing
- No source code access required

**Example:**
```json
{
  "$schema": "./schema/output.schema.json",
  "file_reports": [
    {
      "file": "internal/calculator/math.go",
      "coverage_rate": 0.6667,
      "uncovered_items": [
        {
          "line_range": "15"
        },
        {
          "line_range": "23-27"
        }
      ]
    }
  ]
}
```

See [`docs/examples/minimal.json`](examples/minimal.json) for a complete example.

### Enriched Mode

Combines coverage.out data with AST analysis of source files.

**Characteristics:**
- All UncoveredItem fields populated when available
- Includes `function`, `type`, and `condition` information
- Requires source code access
- Provides deeper insights into uncovered code

**Example:**
```json
{
  "$schema": "./schema/output.schema.json",
  "file_reports": [
    {
      "file": "internal/processor/handler.go",
      "coverage_rate": 0.7250,
      "uncovered_items": [
        {
          "function": "ProcessRequest",
          "type": "branch",
          "condition": "if err != nil",
          "line_range": "34-36"
        },
        {
          "function": "ProcessRequest",
          "type": "error_path",
          "condition": "if status == StatusFailed",
          "line_range": "45"
        }
      ]
    }
  ]
}
```

See [`docs/examples/enriched.json`](examples/enriched.json) for a complete example.

## Validation Rules

### Schema Validation

The schema enforces the following validation rules:

1. **Required fields**: All required fields must be present
2. **Type constraints**: Fields must match their specified types
3. **Range constraints**: `coverage_rate` must be between 0.0 and 1.0
4. **Pattern validation**: `line_range` must match pattern `^[1-9][0-9]*(-[1-9][0-9]*)?$`
5. **Enum validation**: `type` field must be one of the allowed enum values
6. **Additional properties**: No additional properties allowed beyond schema definition

### Line Range Validation

- Start line must be > 0
- End line must be >= start line
- Single-line format: `"45"`
- Multi-line format: `"45-52"`

### Truncation Rules

- Function names: Truncated at 256 characters with `...` suffix
- Conditions: Truncated at 512 characters with `...` suffix

### Sorting Rules

1. **FileReports**: Sorted by `coverage_rate` descending (highest coverage first)
2. **UncoveredItems**: Sorted by line number ascending within each FileReport

## Edge Cases

### 100% Coverage Files

Files with complete coverage are included in the output with an empty `uncovered_items` array:

```json
{
  "file": "internal/config/constants.go",
  "coverage_rate": 1.0,
  "uncovered_items": []
}
```

### 0% Coverage Files

Files with no coverage show all uncovered blocks:

```json
{
  "file": "internal/processor/worker.go",
  "coverage_rate": 0.0,
  "uncovered_items": [
    {
      "line_range": "12"
    },
    {
      "line_range": "18-25"
    }
  ]
}
```

### Files with 0 Statements

Files containing only declarations (no executable statements) have `coverage_rate` set to `1.0`:

```json
{
  "file": "internal/models/types.go",
  "coverage_rate": 1.0,
  "uncovered_items": []
}
```

See [`docs/examples/edge-cases.json`](examples/edge-cases.json) for complete edge case examples.

## Example Files

Three example files demonstrate different aspects of the schema:

1. **[minimal.json](examples/minimal.json)**: Coverage-only mode with basic data
2. **[enriched.json](examples/enriched.json)**: Enriched mode with AST data
3. **[edge-cases.json](examples/edge-cases.json)**: Boundary conditions and special cases

## Schema Validation

The schema can be validated using the `github.com/santhosh-tekuri/jsonschema/v6` library:

```go
import (
    "github.com/santhosh-tekuri/jsonschema/v6"
)

// Load and compile schema
compiler := jsonschema.NewCompiler()
schema, err := compiler.Compile("schema/output.schema.json")
if err != nil {
    // Handle error
}

// Validate JSON data
err = schema.Validate(jsonData)
if err != nil {
    // Handle validation error
}
```

## Future Extensibility

The schema is versioned (`covanalyze-output-schema-v1`) to support future extensions. Future versions may:

- Add new optional fields to existing types
- Introduce new data types
- Extend enum values for the `type` field
- Add metadata fields to the root object

Backward compatibility will be maintained through semantic versioning of the schema `$id`.