package main

import (
	"testing"

	openai "github.com/sashabaranov/go-openai"
)

func TestIsReasoningModel(t *testing.T) {
	tests := []struct {
		model    string
		expected bool
	}{
		{"gpt-4o", false},
		{"gpt-4.1-mini", false},
		{"gpt-5", true},
		{"gpt-5.6-luna", true},
		{"o1-preview", true},
		{"o3-mini", true},
		{"o4-mini", true},
		{"llama3", false},
		{"claude-sonnet-4-5", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			if got := isReasoningModel(tt.model); got != tt.expected {
				t.Errorf("isReasoningModel(%q) = %v, want %v", tt.model, got, tt.expected)
			}
		})
	}
}

func TestBuildChatRequest(t *testing.T) {
	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleUser, Content: "Hello"},
	}

	tests := []struct {
		name                        string
		config                      *Config
		expectedMaxTokens           int
		expectedMaxCompletionTokens int
		expectedTemperature         float32
	}{
		{
			name:                "Standard model uses max_tokens and temperature",
			config:              &Config{Model: "gpt-4o", MaxTokens: 1000, Temperature: 0.7},
			expectedMaxTokens:   1000,
			expectedTemperature: 0.7,
		},
		{
			name: "Reasoning model remaps max_tokens",
			config: &Config{
				Model:       "gpt-5.6-luna",
				MaxTokens:   1000,
				Temperature: 0.7,
			},
			expectedMaxCompletionTokens: 1000,
			expectedTemperature:         0,
		},
		{
			name: "Explicit max_completion_tokens takes precedence",
			config: &Config{
				Model:               "gpt-4o",
				MaxTokens:           1000,
				MaxCompletionTokens: 2000,
				Temperature:         0.7,
			},
			expectedMaxCompletionTokens: 2000,
			expectedTemperature:         0.7,
		},
		{
			name: "Reasoning model with explicit max_completion_tokens",
			config: &Config{
				Model:               "o3-mini",
				MaxTokens:           1000,
				MaxCompletionTokens: 2000,
				Temperature:         1,
			},
			expectedMaxCompletionTokens: 2000,
			expectedTemperature:         0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := buildChatRequest(tt.config, messages, nil)

			if req.MaxTokens != tt.expectedMaxTokens {
				t.Errorf("MaxTokens = %d, want %d", req.MaxTokens, tt.expectedMaxTokens)
			}
			if req.MaxCompletionTokens != tt.expectedMaxCompletionTokens {
				t.Errorf(
					"MaxCompletionTokens = %d, want %d",
					req.MaxCompletionTokens,
					tt.expectedMaxCompletionTokens,
				)
			}
			if req.Temperature != tt.expectedTemperature {
				t.Errorf("Temperature = %f, want %f", req.Temperature, tt.expectedTemperature)
			}
			if req.Model != tt.config.Model {
				t.Errorf("Model = %q, want %q", req.Model, tt.config.Model)
			}
		})
	}
}

func TestBuildChatRequestWithTool(t *testing.T) {
	messages := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleUser, Content: "Hello"},
	}
	toolMeta, err := ParseToolSchema(
		`{"name": "test_tool", "parameters": {"type": "object", "properties": {"result": {"type": "string"}}}}`,
	)
	if err != nil {
		t.Fatalf("failed to parse tool schema: %v", err)
	}

	config := &Config{Model: "gpt-4o", MaxTokens: 1000, Temperature: 0.7}
	req := buildChatRequest(config, messages, toolMeta)

	if len(req.Tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(req.Tools))
	}
	if req.ToolChoice == nil {
		t.Fatal("expected ToolChoice to be set")
	}
}
