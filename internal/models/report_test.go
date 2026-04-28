package models

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// TestSchemaValidation tests that the JSON Schema loads and compiles successfully
func TestSchemaValidation(t *testing.T) {
	schemaPath := filepath.Join("..", "..", "schema", "output.schema.json")

	// Load schema file
	schemaData, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("Failed to read schema file: %v", err)
	}

	// Parse schema JSON
	var schemaObj interface{}
	if err := json.Unmarshal(schemaData, &schemaObj); err != nil {
		t.Fatalf("Failed to parse schema JSON: %v", err)
	}

	// Compile schema
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", schemaObj); err != nil {
		t.Fatalf("Failed to add schema resource: %v", err)
	}

	schema, err := compiler.Compile("schema.json")
	if err != nil {
		t.Fatalf("Failed to compile schema: %v", err)
	}

	if schema == nil {
		t.Fatal("Schema is nil after compilation")
	}
}

// TestValidateMinimalExample tests that minimal.json validates against the schema
func TestValidateMinimalExample(t *testing.T) {
	examplePath := filepath.Join("..", "..", "docs", "examples", "minimal.json")
	validateExampleFile(t, examplePath)
}

// TestValidateEnrichedExample tests that enriched.json validates against the schema
func TestValidateEnrichedExample(t *testing.T) {
	examplePath := filepath.Join("..", "..", "docs", "examples", "enriched.json")
	validateExampleFile(t, examplePath)
}

// TestValidateEdgeCasesExample tests that edge-cases.json validates against the schema
func TestValidateEdgeCasesExample(t *testing.T) {
	examplePath := filepath.Join("..", "..", "docs", "examples", "edge-cases.json")
	validateExampleFile(t, examplePath)
}

// validateExampleFile is a helper function to validate an example file against the schema
func validateExampleFile(t *testing.T, examplePath string) {
	t.Helper()

	// Load schema
	schemaPath := filepath.Join("..", "..", "schema", "output.schema.json")
	schemaData, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("Failed to read schema file: %v", err)
	}

	// Parse schema JSON
	var schemaObj interface{}
	if err := json.Unmarshal(schemaData, &schemaObj); err != nil {
		t.Fatalf("Failed to parse schema JSON: %v", err)
	}

	// Compile schema
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", schemaObj); err != nil {
		t.Fatalf("Failed to add schema resource: %v", err)
	}

	schema, err := compiler.Compile("schema.json")
	if err != nil {
		t.Fatalf("Failed to compile schema: %v", err)
	}

	// Load example file
	exampleData, err := os.ReadFile(examplePath)
	if err != nil {
		t.Fatalf("Failed to read example file %s: %v", examplePath, err)
	}

	// Parse JSON
	var data interface{}
	if err := json.Unmarshal(exampleData, &data); err != nil {
		t.Fatalf("Failed to parse example JSON from %s: %v", examplePath, err)
	}

	// Validate against schema
	if err := schema.Validate(data); err != nil {
		t.Errorf("Example file %s failed validation: %v", filepath.Base(examplePath), err)
	}
}

// TestInvalidDataValidation tests that invalid data is correctly rejected by the schema
func TestInvalidDataValidation(t *testing.T) {
	// Load schema
	schemaPath := filepath.Join("..", "..", "schema", "output.schema.json")
	schemaData, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("Failed to read schema file: %v", err)
	}

	// Parse schema JSON
	var schemaObj interface{}
	if err := json.Unmarshal(schemaData, &schemaObj); err != nil {
		t.Fatalf("Failed to parse schema JSON: %v", err)
	}

	// Compile schema
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", schemaObj); err != nil {
		t.Fatalf("Failed to add schema resource: %v", err)
	}

	schema, err := compiler.Compile("schema.json")
	if err != nil {
		t.Fatalf("Failed to compile schema: %v", err)
	}

	tests := []struct {
		name        string
		data        string
		shouldFail  bool
		description string
	}{
		{
			name: "missing_line_range",
			data: `{
				"$schema": "./schema/output.schema.json",
				"file_reports": [{
					"file": "test.go",
					"coverage_rate": 0.5,
					"uncovered_items": [{
						"function": "TestFunc"
					}]
				}]
			}`,
			shouldFail:  true,
			description: "UncoveredItem missing required line_range field",
		},
		{
			name: "invalid_coverage_rate_above_1",
			data: `{
				"$schema": "./schema/output.schema.json",
				"file_reports": [{
					"file": "test.go",
					"coverage_rate": 1.5,
					"uncovered_items": []
				}]
			}`,
			shouldFail:  true,
			description: "coverage_rate above maximum value of 1.0",
		},
		{
			name: "invalid_coverage_rate_negative",
			data: `{
				"$schema": "./schema/output.schema.json",
				"file_reports": [{
					"file": "test.go",
					"coverage_rate": -0.1,
					"uncovered_items": []
				}]
			}`,
			shouldFail:  true,
			description: "coverage_rate below minimum value of 0.0",
		},
		{
			name: "invalid_type_enum",
			data: `{
				"$schema": "./schema/output.schema.json",
				"file_reports": [{
					"file": "test.go",
					"coverage_rate": 0.5,
					"uncovered_items": [{
						"line_range": "10",
						"type": "invalid_type"
					}]
				}]
			}`,
			shouldFail:  true,
			description: "type field with invalid enum value",
		},
		{
			name: "invalid_line_range_format",
			data: `{
				"$schema": "./schema/output.schema.json",
				"file_reports": [{
					"file": "test.go",
					"coverage_rate": 0.5,
					"uncovered_items": [{
						"line_range": "0"
					}]
				}]
			}`,
			shouldFail:  true,
			description: "line_range with invalid format (line 0)",
		},
		{
			name: "valid_minimal_data",
			data: `{
				"$schema": "./schema/output.schema.json",
				"file_reports": [{
					"file": "test.go",
					"coverage_rate": 0.5000,
					"uncovered_items": [{
						"line_range": "10"
					}]
				}]
			}`,
			shouldFail:  false,
			description: "Valid minimal data should pass validation",
		},
		{
			name: "valid_enriched_data",
			data: `{
				"$schema": "./schema/output.schema.json",
				"file_reports": [{
					"file": "test.go",
					"coverage_rate": 0.7500,
					"uncovered_items": [{
						"function": "TestFunc",
						"type": "branch",
						"condition": "if x > 0",
						"line_range": "10-15"
					}]
				}]
			}`,
			shouldFail:  false,
			description: "Valid enriched data with all optional fields should pass validation",
		},
		{
			name: "valid_empty_type",
			data: `{
				"$schema": "./schema/output.schema.json",
				"file_reports": [{
					"file": "test.go",
					"coverage_rate": 0.5000,
					"uncovered_items": [{
						"type": "",
						"line_range": "10"
					}]
				}]
			}`,
			shouldFail:  false,
			description: "Empty string for type field should be valid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data interface{}
			if err := json.Unmarshal([]byte(tt.data), &data); err != nil {
				t.Fatalf("Failed to parse test JSON: %v", err)
			}

			err := schema.Validate(data)
			if tt.shouldFail && err == nil {
				t.Errorf("Expected validation to fail for %s, but it passed", tt.description)
			}
			if !tt.shouldFail && err != nil {
				t.Errorf("Expected validation to pass for %s, but it failed: %v", tt.description, err)
			}
		})
	}
}
