package main

import (
	"testing"

	openai "github.com/sashabaranov/go-openai"
)

func TestParseToolSchema(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		validate    func(*testing.T, *ToolMeta)
	}{
		{
			name: "Valid tool schema",
			input: `{
				"name": "get_city_info",
				"description": "Get information about a city",
				"parameters": {
					"type": "object",
					"properties": {
						"city": { "type": "string" }
					},
					"required": ["city"]
				}
			}`,
			expectError: false,
			validate: func(t *testing.T, meta *ToolMeta) {
				if meta.Name != "get_city_info" {
					t.Errorf("expected name 'get_city_info', got '%s'", meta.Name)
				}
				if meta.Description != "Get information about a city" {
					t.Errorf(
						"expected description 'Get information about a city', got '%s'",
						meta.Description,
					)
				}
				if meta.Parameters == nil {
					t.Error("expected parameters to be non-nil")
				}
				if meta.Parameters["type"] != "object" {
					t.Errorf("expected parameters type 'object', got '%v'", meta.Parameters["type"])
				}
			},
		},
		{
			name: "Minimal tool schema",
			input: `{
				"name": "simple_function",
				"parameters": {}
			}`,
			expectError: false,
			validate: func(t *testing.T, meta *ToolMeta) {
				if meta.Name != "simple_function" {
					t.Errorf("expected name 'simple_function', got '%s'", meta.Name)
				}
				if meta.Description != "" {
					t.Errorf("expected empty description, got '%s'", meta.Description)
				}
			},
		},
		{
			name:        "Empty string",
			input:       "",
			expectError: false,
			validate: func(t *testing.T, meta *ToolMeta) {
				if meta != nil {
					t.Error("expected nil for empty input")
				}
			},
		},
		{
			name:        "Invalid JSON",
			input:       "not valid json",
			expectError: true,
		},
		{
			name:        "Missing name field",
			input:       `{"description": "test", "parameters": {}}`,
			expectError: true,
		},
		{
			name:        "Empty name field",
			input:       `{"name": "", "parameters": {}}`,
			expectError: true,
		},
		{
			name:        "Malformed JSON",
			input:       `{"name": "test"`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta, err := ParseToolSchema(tt.input)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.expectError && tt.validate != nil {
				tt.validate(t, meta)
			}
		})
	}
}

func TestToOpenAITool(t *testing.T) {
	meta := &ToolMeta{
		Name:        "get_weather",
		Description: "Get the current weather for a location",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"location": map[string]interface{}{
					"type":        "string",
					"description": "The city and state, e.g. San Francisco, CA",
				},
				"unit": map[string]interface{}{
					"type": "string",
					"enum": []string{"celsius", "fahrenheit"},
				},
			},
			"required": []string{"location"},
		},
	}

	tool := meta.ToOpenAITool()

	if tool.Type != openai.ToolTypeFunction {
		t.Errorf("expected tool type '%s', got '%s'", openai.ToolTypeFunction, tool.Type)
	}

	if tool.Function == nil {
		t.Fatal("expected function to be non-nil")
	}

	if tool.Function.Name != "get_weather" {
		t.Errorf("expected function name 'get_weather', got '%s'", tool.Function.Name)
	}

	if tool.Function.Description != "Get the current weather for a location" {
		t.Errorf(
			"expected description 'Get the current weather for a location', got '%s'",
			tool.Function.Description,
		)
	}

	if tool.Function.Parameters == nil {
		t.Error("expected parameters to be non-nil")
	}
}

func TestToOpenAIToolMinimal(t *testing.T) {
	meta := &ToolMeta{
		Name:       "minimal_function",
		Parameters: map[string]interface{}{},
	}

	tool := meta.ToOpenAITool()

	if tool.Type != openai.ToolTypeFunction {
		t.Errorf("expected tool type '%s', got '%s'", openai.ToolTypeFunction, tool.Type)
	}

	if tool.Function.Name != "minimal_function" {
		t.Errorf("expected function name 'minimal_function', got '%s'", tool.Function.Name)
	}

	if tool.Function.Description != "" {
		t.Errorf("expected empty description, got '%s'", tool.Function.Description)
	}
}

func TestParseToolSchemaWithComplexParameters(t *testing.T) {
	input := `{
		"name": "analyze_code",
		"description": "Analyze code and return structured results",
		"parameters": {
			"type": "object",
			"properties": {
				"score": {
					"type": "integer",
					"description": "Code quality score 1-10"
				},
				"issues": {
					"type": "array",
					"items": { "type": "string" },
					"description": "List of issues found"
				},
				"metadata": {
					"type": "object",
					"properties": {
						"language": { "type": "string" },
						"lines": { "type": "integer" }
					}
				},
				"is_valid": {
					"type": "boolean"
				}
			},
			"required": ["score", "issues"]
		}
	}`

	meta, err := ParseToolSchema(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if meta.Name != "analyze_code" {
		t.Errorf("expected name 'analyze_code', got '%s'", meta.Name)
	}

	// Verify parameters structure
	params := meta.Parameters
	if params == nil {
		t.Fatal("expected parameters to be non-nil")
	}

	props, ok := params["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("expected properties to be a map")
	}

	// Check score property
	score, ok := props["score"].(map[string]interface{})
	if !ok {
		t.Fatal("expected score property to be a map")
	}
	if score["type"] != "integer" {
		t.Errorf("expected score type 'integer', got '%v'", score["type"])
	}

	// Check issues property (array)
	issues, ok := props["issues"].(map[string]interface{})
	if !ok {
		t.Fatal("expected issues property to be a map")
	}
	if issues["type"] != "array" {
		t.Errorf("expected issues type 'array', got '%v'", issues["type"])
	}
}

// TestParseFunctionArguments tests parsing of function call arguments JSON
func TestParseFunctionArguments(t *testing.T) {
	tests := []struct {
		name        string
		jsonStr     string
		expected    map[string]string
		expectError bool
	}{
		{
			name:    "Simple string values",
			jsonStr: `{"city": "Paris", "country": "France"}`,
			expected: map[string]string{
				"city":    "Paris",
				"country": "France",
			},
		},
		{
			name:    "Integer value",
			jsonStr: `{"score": 8, "name": "test"}`,
			expected: map[string]string{
				"score": "8",
				"name":  "test",
			},
		},
		{
			name:    "Boolean value",
			jsonStr: `{"is_valid": true, "is_error": false}`,
			expected: map[string]string{
				"is_valid": "true",
				"is_error": "false",
			},
		},
		{
			name:    "Array value",
			jsonStr: `{"issues": ["bug1", "bug2"], "name": "review"}`,
			expected: map[string]string{
				"issues": `["bug1","bug2"]`,
				"name":   "review",
			},
		},
		{
			name:    "Nested object value",
			jsonStr: `{"metadata": {"lang": "go", "version": 1}, "name": "test"}`,
			expected: map[string]string{
				"metadata": `{"lang":"go","version":1}`,
				"name":     "test",
			},
		},
		{
			name:    "Float value",
			jsonStr: `{"temperature": 0.7, "name": "config"}`,
			expected: map[string]string{
				"temperature": "0.7",
				"name":        "config",
			},
		},
		{
			name:    "Null value",
			jsonStr: `{"optional": null, "name": "test"}`,
			expected: map[string]string{
				"optional": "null",
				"name":     "test",
			},
		},
		{
			name:     "Empty object",
			jsonStr:  `{}`,
			expected: map[string]string{},
		},
		{
			name:     "Empty string",
			jsonStr:  "",
			expected: map[string]string{},
		},
		{
			name:        "Invalid JSON",
			jsonStr:     "not valid json",
			expectError: true,
		},
		{
			name:        "Malformed JSON",
			jsonStr:     `{"key": "value"`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := ParseFunctionArguments(tt.jsonStr)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify each expected key-value pair
			for key, expectedValue := range tt.expected {
				actualValue, ok := output[key]
				if !ok {
					t.Errorf("expected key '%s' not found in output", key)
					continue
				}
				if actualValue != expectedValue {
					t.Errorf(
						"for key '%s': expected '%s', got '%s'",
						key,
						expectedValue,
						actualValue,
					)
				}
			}

			// Verify no extra keys
			if len(output) != len(tt.expected) {
				t.Errorf("expected %d keys, got %d", len(tt.expected), len(output))
			}
		})
	}
}

func TestParseToolSchemaWithUnicodeCharacters(t *testing.T) {
	input := `{
		"name": "get_info",
		"description": "取得資訊 - 获取信息",
		"parameters": {
			"type": "object",
			"properties": {
				"城市": { "type": "string", "description": "城市名稱" }
			}
		}
	}`

	meta, err := ParseToolSchema(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if meta.Name != "get_info" {
		t.Errorf("expected name 'get_info', got '%s'", meta.Name)
	}

	if meta.Description != "取得資訊 - 获取信息" {
		t.Errorf("expected description with unicode, got '%s'", meta.Description)
	}
}

func TestParseToolSchemaWithSpecialCharacters(t *testing.T) {
	input := `{
		"name": "handle_text",
		"description": "Handle text with special chars: \"quotes\", \\backslash, \n newline",
		"parameters": {
			"type": "object",
			"properties": {}
		}
	}`

	meta, err := ParseToolSchema(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if meta.Name != "handle_text" {
		t.Errorf("expected name 'handle_text', got '%s'", meta.Name)
	}

	expectedDesc := "Handle text with special chars: \"quotes\", \\backslash, \n newline"
	if meta.Description != expectedDesc {
		t.Errorf("expected description '%s', got '%s'", expectedDesc, meta.Description)
	}
}
