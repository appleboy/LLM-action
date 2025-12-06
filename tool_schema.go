package main

import (
	"encoding/json"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

// ToolMeta represents the function schema structure for structured output
type ToolMeta struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// ParseToolSchema parses JSON string to ToolMeta
func ParseToolSchema(jsonStr string) (*ToolMeta, error) {
	if jsonStr == "" {
		return nil, nil
	}

	var meta ToolMeta
	if err := json.Unmarshal([]byte(jsonStr), &meta); err != nil {
		return nil, fmt.Errorf("failed to parse tool_schema JSON: %w", err)
	}

	if meta.Name == "" {
		return nil, fmt.Errorf("tool_schema must have a 'name' field")
	}

	return &meta, nil
}

// ToOpenAITool converts ToolMeta to openai.Tool format
func (t *ToolMeta) ToOpenAITool() openai.Tool {
	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        t.Name,
			Description: t.Description,
			Parameters:  t.Parameters,
		},
	}
}

// ParseFunctionArguments parses function call arguments JSON string
// and converts it to a map[string]string for GitHub Actions output.
// String values are kept as-is, other types are marshaled back to JSON.
func ParseFunctionArguments(jsonStr string) (map[string]string, error) {
	if jsonStr == "" {
		return map[string]string{}, nil
	}

	var args map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &args); err != nil {
		return nil, fmt.Errorf("failed to parse function arguments: %w", err)
	}

	output := make(map[string]string)
	for key, value := range args {
		switch v := value.(type) {
		case string:
			output[key] = v
		default:
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal value for key '%s': %w", key, err)
			}
			output[key] = string(jsonBytes)
		}
	}

	return output, nil
}
