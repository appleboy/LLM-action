package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

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
	const maskPattern = "********"

	if len(apiKey) <= 8 {
		return maskPattern
	}
	// Show first 4 and last 4 characters with fixed mask pattern
	return apiKey[:4] + maskPattern + apiKey[len(apiKey)-4:]
}

// printTokenUsage prints token usage statistics to stdout
func printTokenUsage(usage openai.Usage) {
	fmt.Println("--- Token Usage ---")
	fmt.Printf("Prompt Tokens: %d\n", usage.PromptTokens)
	fmt.Printf("Completion Tokens: %d\n", usage.CompletionTokens)
	fmt.Printf("Total Tokens: %d\n", usage.TotalTokens)
	if usage.PromptTokensDetails != nil {
		fmt.Printf("Cached Tokens: %d\n", usage.PromptTokensDetails.CachedTokens)
	}
	if d := usage.CompletionTokensDetails; d != nil {
		fmt.Printf("Reasoning Tokens: %d\n", d.ReasoningTokens)
		fmt.Printf("Accepted Prediction Tokens: %d\n", d.AcceptedPredictionTokens)
		fmt.Printf("Rejected Prediction Tokens: %d\n", d.RejectedPredictionTokens)
	}
	fmt.Println("--- End Token Usage ---")
}

// addTokenUsageToOutput adds token usage metrics to the output map
func addTokenUsageToOutput(output map[string]string, usage openai.Usage) {
	output["prompt_tokens"] = strconv.Itoa(usage.PromptTokens)
	output["completion_tokens"] = strconv.Itoa(usage.CompletionTokens)
	output["total_tokens"] = strconv.Itoa(usage.TotalTokens)

	if usage.PromptTokensDetails != nil {
		output["prompt_cached_tokens"] = strconv.Itoa(usage.PromptTokensDetails.CachedTokens)
	}

	if d := usage.CompletionTokensDetails; d != nil {
		output["completion_reasoning_tokens"] = strconv.Itoa(d.ReasoningTokens)
		output["completion_accepted_prediction_tokens"] = strconv.Itoa(d.AcceptedPredictionTokens)
		output["completion_rejected_prediction_tokens"] = strconv.Itoa(d.RejectedPredictionTokens)
	}
}

// extractResponse extracts the response content from the API response
func extractResponse(
	resp openai.ChatCompletionResponse,
	toolMeta *ToolMeta,
	debug bool,
) (string, error) {
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	if toolMeta != nil {
		// Extract function call arguments when tool schema is used
		if len(resp.Choices[0].Message.ToolCalls) > 0 {
			// Debug: Print tool call details if debug mode is enabled
			if debug {
				fmt.Println("=== Debug Mode: Tool Calls ===")
				if err := godump.Dump(resp.Choices[0].Message.ToolCalls); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to dump tool calls: %v\n", err)
				}
				fmt.Println("==============================")
			}
			return resp.Choices[0].Message.ToolCalls[0].Function.Arguments, nil
		}
		return "", fmt.Errorf("expected tool call response but got none")
	}

	return resp.Choices[0].Message.Content, nil
}

// prepareToolSchema parses and validates the tool schema if provided
func prepareToolSchema(config *Config) (*ToolMeta, error) {
	if config.ToolSchema == "" {
		return nil, nil
	}

	toolMeta, err := ParseToolSchema(config.ToolSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tool schema: %w", err)
	}

	// Debug: Print tool schema if debug mode is enabled
	if config.Debug {
		fmt.Println("=== Debug Mode: Tool Schema ===")
		if err := godump.Dump(toolMeta); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to dump tool schema: %v\n", err)
		}
		fmt.Println("===============================")
	}

	return toolMeta, nil
}

// reasoningModelPrefixes mirrors go-openai's ReasoningValidator: these model
// families reject max_tokens and non-default temperature before any request
// is sent.
var reasoningModelPrefixes = []string{"o1", "o3", "o4", "gpt-5"}

// isReasoningModel reports whether the model only accepts reasoning-model
// parameters (max_completion_tokens instead of max_tokens, fixed temperature).
func isReasoningModel(model string) bool {
	for _, prefix := range reasoningModelPrefixes {
		if strings.HasPrefix(model, prefix) {
			return true
		}
	}
	return false
}

// buildChatRequest creates a chat completion request with optional tool support
func buildChatRequest(
	config *Config,
	messages []openai.ChatCompletionMessage,
	toolMeta *ToolMeta,
) openai.ChatCompletionRequest {
	req := openai.ChatCompletionRequest{
		Model:    config.Model,
		Messages: messages,
	}

	switch {
	case config.MaxCompletionTokens > 0:
		req.MaxCompletionTokens = config.MaxCompletionTokens
	case isReasoningModel(config.Model):
		// Reasoning models reject max_tokens — send the same budget through
		// max_completion_tokens so existing configs keep working.
		req.MaxCompletionTokens = config.MaxTokens
	default:
		req.MaxTokens = config.MaxTokens
	}

	if isReasoningModel(config.Model) {
		// Temperature is fixed at 1 for reasoning models, sending any other
		// value fails the request.
		if config.Temperature != 1 {
			fmt.Println("Note: temperature is fixed at 1 for reasoning models, omitting it")
		}
	} else {
		req.Temperature = float32(config.Temperature)
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

	return req
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
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Build messages
	messages := BuildMessages(config)

	// Parse and validate tool schema if provided
	toolMeta, err := prepareToolSchema(config)
	if err != nil {
		return err
	}

	// Debug: Print messages if debug mode is enabled
	if config.Debug {
		fmt.Println("=== Debug Mode: Messages ===")
		if err := godump.Dump(messages); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to dump messages: %v\n", err)
		}
		fmt.Println("============================")
	}

	// Create chat completion request with optional tool support
	req := buildChatRequest(config, messages, toolMeta)

	fmt.Println("Sending request to LLM...")
	fmt.Printf("Model: %s\n", config.Model)
	fmt.Printf("Base URL: %s\n", config.BaseURL)

	// Call the API
	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return fmt.Errorf("chat completion error: %w", err)
	}

	// Extract response content
	response, err := extractResponse(resp, toolMeta, config.Debug)
	if err != nil {
		return err
	}

	// Print response for debugging
	fmt.Println("--- LLM Response ---")
	fmt.Println(response)
	fmt.Println("--- End Response ---")

	// Print token usage statistics
	printTokenUsage(resp.Usage)

	// Set GitHub Actions output
	var toolArgs map[string]string
	if toolMeta != nil {
		// Parse JSON arguments
		var err error
		toolArgs, err = ParseFunctionArguments(response)
		if err != nil {
			return fmt.Errorf("failed to parse function arguments: %w", err)
		}
	}

	// Build output map with raw response and tool arguments
	output, reservedFieldSkipped := BuildOutputMap(response, toolArgs)
	if reservedFieldSkipped {
		fmt.Fprintf(
			os.Stderr,
			"Warning: tool schema field '%s' is reserved and will be skipped\n",
			ReservedOutputField,
		)
	}

	// Add token usage metrics to output
	addTokenUsageToOutput(output, resp.Usage)

	if err := gh.SetOutput(output); err != nil {
		return fmt.Errorf("failed to set output: %w", err)
	}

	return nil
}
