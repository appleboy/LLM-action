package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

// RenderTemplate renders a Go template string with environment variables as data
// Environment variables with INPUT_ prefix are available both with and without the prefix
// For example: INPUT_MODEL can be accessed as {{.MODEL}} or {{.INPUT_MODEL}}
func RenderTemplate(templateStr string) (string, error) {
	// Build template data from environment variables
	data := buildTemplateData()

	// Parse template
	tmpl, err := template.New("prompt").Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// buildTemplateData builds a map of environment variables for template rendering
// INPUT_ prefixed variables are available both with and without the prefix
func buildTemplateData() map[string]string {
	data := make(map[string]string)

	// Get all environment variables
	environ := os.Environ()
	for _, env := range environ {
		// Split key=value
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		// Add with original key
		data[key] = value

		// If key starts with INPUT_, also add without prefix
		if strings.HasPrefix(key, "INPUT_") {
			keyWithoutPrefix := strings.TrimPrefix(key, "INPUT_")
			data[keyWithoutPrefix] = value
		}
	}

	return data
}
