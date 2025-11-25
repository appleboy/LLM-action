package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

var (
	errAPIKeyRequired      = errors.New("api_key is required")
	errInputPromptRequired = errors.New("input_prompt is required")
)

// Config holds all configuration for the LLM action
type Config struct {
	BaseURL       string
	APIKey        string
	Model         string
	SkipSSLVerify bool
	SystemPrompt  string
	InputPrompt   string
	Temperature   float64
	MaxTokens     int
	Debug         bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		BaseURL:     os.Getenv("INPUT_BASE_URL"),
		APIKey:      os.Getenv("INPUT_API_KEY"),
		Model:       os.Getenv("INPUT_MODEL"),
		Temperature: 0.7,  // default
		MaxTokens:   1000, // default
	}

	// Set default base URL if not provided
	if config.BaseURL == "" {
		config.BaseURL = "https://api.openai.com/v1"
	}

	// Validate required inputs
	if config.APIKey == "" {
		return nil, errAPIKeyRequired
	}

	// Load input prompt (supports text, file path, or URL)
	inputPromptInput := os.Getenv("INPUT_INPUT_PROMPT")
	if inputPromptInput == "" {
		return nil, errInputPromptRequired
	}
	loadedInputPrompt, err := LoadPrompt(inputPromptInput)
	if err != nil {
		return nil, fmt.Errorf("failed to load input_prompt: %w", err)
	}
	config.InputPrompt = loadedInputPrompt

	// Load system prompt (supports text, file path, or URL)
	systemPromptInput := os.Getenv("INPUT_SYSTEM_PROMPT")
	if systemPromptInput != "" {
		loadedPrompt, err := LoadPrompt(systemPromptInput)
		if err != nil {
			return nil, fmt.Errorf("failed to load system_prompt: %w", err)
		}
		config.SystemPrompt = loadedPrompt
	}

	// Parse optional parameters
	if err := config.parseTemperature(os.Getenv("INPUT_TEMPERATURE")); err != nil {
		return nil, err
	}

	if err := config.parseMaxTokens(os.Getenv("INPUT_MAX_TOKENS")); err != nil {
		return nil, err
	}

	if err := config.parseSkipSSL(os.Getenv("INPUT_SKIP_SSL_VERIFY")); err != nil {
		return nil, err
	}

	if err := config.parseDebug(os.Getenv("INPUT_DEBUG")); err != nil {
		return nil, err
	}

	return config, nil
}

// parseTemperature parses temperature string to float64
func (c *Config) parseTemperature(s string) error {
	if s == "" {
		return nil
	}

	temp, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return fmt.Errorf("invalid temperature value: %v", err)
	}
	c.Temperature = temp
	return nil
}

// parseMaxTokens parses max tokens string to int
func (c *Config) parseMaxTokens(s string) error {
	if s == "" {
		return nil
	}

	tokens, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("invalid max_tokens value: %v", err)
	}
	if tokens < 0 {
		return fmt.Errorf("max_tokens must be positive")
	}
	c.MaxTokens = tokens
	return nil
}

// parseSkipSSL parses skip SSL verify string to bool
func (c *Config) parseSkipSSL(s string) error {
	if s == "" {
		return nil
	}

	skip, err := strconv.ParseBool(s)
	if err != nil {
		return fmt.Errorf("invalid skip_ssl_verify value: %v", err)
	}
	c.SkipSSLVerify = skip
	return nil
}

// parseDebug parses debug string to bool
func (c *Config) parseDebug(s string) error {
	if s == "" {
		return nil
	}

	debug, err := strconv.ParseBool(s)
	if err != nil {
		return fmt.Errorf("invalid debug value: %v", err)
	}
	c.Debug = debug
	return nil
}
