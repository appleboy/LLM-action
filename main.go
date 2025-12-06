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

	// Debug: Print messages if debug mode is enabled
	if config.Debug {
		fmt.Println("=== Debug Mode: Messages ===")
		if err := godump.Dump(messages); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to dump messages: %v\n", err)
		}
		fmt.Println("============================")
	}

	// Create chat completion request
	req := openai.ChatCompletionRequest{
		Model:       config.Model,
		Messages:    messages,
		Temperature: float32(config.Temperature),
		MaxTokens:   config.MaxTokens,
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

	response := resp.Choices[0].Message.Content

	// Print response for debugging
	fmt.Println("--- LLM Response ---")
	fmt.Println(response)
	fmt.Println("--- End Response ---")

	// Set GitHub Actions output
	if err := gh.SetOutput(map[string]string{
		"response": response,
	}); err != nil {
		return fmt.Errorf("failed to set output: %v", err)
	}

	return nil
}
