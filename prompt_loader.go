package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// LoadPrompt intelligently loads prompt content from text, file, or URL
// It detects the input type automatically:
// - If starts with http:// or https:// -> loads from URL
// - If starts with file:// or is a valid file path -> loads from file
// - Otherwise -> returns as plain text
// After loading, it renders the content as a Go template with environment variables
func LoadPrompt(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	var content string
	var err error

	// Determine source type and load content
	switch {
	case isURL(input):
		content, err = loadFromURL(input)
		if err != nil {
			return "", err
		}
	case isFilePath(input):
		content, err = loadFromFile(input)
		if err != nil {
			return "", err
		}
	default:
		// Return as plain text
		content = input
	}

	// Render template with environment variables
	rendered, err := RenderTemplate(content)
	if err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return rendered, nil
}

// isURL checks if the input string is a URL
func isURL(input string) bool {
	return strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://")
}

// isFilePath checks if the input string is a file path
func isFilePath(input string) bool {
	// Remove file:// prefix if present
	path := strings.TrimPrefix(input, "file://")

	// Check if file exists
	if _, err := os.Stat(path); err == nil {
		return true
	}

	// If it starts with file://, treat it as a file path even if it doesn't exist
	// This will allow proper error reporting
	return strings.HasPrefix(input, "file://")
}

// loadFromURL loads content from a URL
func loadFromURL(url string) (string, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create request with context
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request for URL %s: %w", url, err)
	}

	// Add User-Agent header
	req.Header.Set("User-Agent", "LLM-Action/1.0")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL %s: %w", url, err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch URL %s: status code %d", url, resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response from URL %s: %w", url, err)
	}

	return string(body), nil
}

// loadFromFile loads content from a local file
func loadFromFile(path string) (string, error) {
	// Remove file:// prefix if present
	cleanPath := strings.TrimPrefix(path, "file://")

	// Read file
	content, err := os.ReadFile(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", cleanPath, err)
	}

	return string(content), nil
}
