package main

import (
	"context"
	"fmt"
	"os"

	"github.com/appleboy/com/gh"
	openai "github.com/sashabaranov/go-openai"
	"github.com/yassinebenaid/godump"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// maskAPIKey masks the API key for secure logging
func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return "********"
	}
	// Show first 4 and last 4 characters
	return apiKey[:4] + "****" + apiKey[len(apiKey)-4:]
}

func run() error {
	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	// Debug: Print all parameters if debug mode is enabled
	if config.Debug {
		fmt.Println("=== Debug Mode: All Parameters ===")
		// Create a copy of config with masked API key for secure logging
		debugConfig := *config
		debugConfig.APIKey = maskAPIKey(config.APIKey)
		if err := godump.Dump(debugConfig); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to dump config: %v\n", err)
		}
		fmt.Println("===================================")
	}

	// Create OpenAI client
	client, err := NewClient(config)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}

	// Build messages
	messages := BuildMessages(config)

	// Parse tool schema if provided
	var toolMeta *ToolMeta
	if config.ToolSchema != "" {
		toolMeta, err = ParseToolSchema(config.ToolSchema)
		if err != nil {
			return fmt.Errorf("failed to parse tool schema: %v", err)
		}
	}

	// Debug: Print messages if debug mode is enabled
	if config.Debug {
		fmt.Println("=== Debug Mode: Messages ===")
		if err := godump.Dump(messages); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to dump messages: %v\n", err)
		}
		fmt.Println("============================")

		if toolMeta != nil {
			fmt.Println("=== Debug Mode: Tool Schema ===")
			if err := godump.Dump(toolMeta); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to dump tool schema: %v\n", err)
			}
			fmt.Println("===============================")
		}
	}

	// Create chat completion request
	req := openai.ChatCompletionRequest{
		Model:       config.Model,
		Messages:    messages,
		Temperature: float32(config.Temperature),
		MaxTokens:   config.MaxTokens,
	}

	// Add tool if schema provided
	if toolMeta != nil {
		req.Tools = []openai.Tool{toolMeta.ToOpenAITool()}
		// Force the model to use this specific function
		req.ToolChoice = &openai.ToolChoice{
			Type: openai.ToolTypeFunction,
			Function: openai.ToolFunction{
				Name: toolMeta.Name,
			},
		}
	}

	fmt.Println("Sending request to LLM...")
	fmt.Printf("Model: %s\n", config.Model)
	fmt.Printf("Base URL: %s\n", config.BaseURL)

	// Call the API
	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return fmt.Errorf("chat completion error: %v", err)
	}

	// Extract response
	if len(resp.Choices) == 0 {
		return fmt.Errorf("no response from LLM")
	}

	var response string
	if toolMeta != nil {
		// Extract function call arguments when tool schema is used
		if len(resp.Choices[0].Message.ToolCalls) > 0 {
			response = resp.Choices[0].Message.ToolCalls[0].Function.Arguments
		} else {
			return fmt.Errorf("expected tool call response but got none")
		}
	} else {
		response = resp.Choices[0].Message.Content
	}

	// Print response for debugging
	fmt.Println("--- LLM Response ---")
	fmt.Println(response)
	fmt.Println("--- End Response ---")

	// Set GitHub Actions output
	var output map[string]string
	if toolMeta != nil {
		// Parse JSON arguments and set each field as output
		var err error
		output, err = ParseFunctionArguments(response)
		if err != nil {
			return err
		}
	} else {
		output = map[string]string{
			"response": response,
		}
	}

	if err := gh.SetOutput(output); err != nil {
		return fmt.Errorf("failed to set output: %v", err)
	}

	return nil
}
