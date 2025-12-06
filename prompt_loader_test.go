package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadPrompt_PlainText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple text",
			input:    "You are a helpful assistant",
			expected: "You are a helpful assistant",
		},
		{
			name:     "multiline text",
			input:    "Line 1\nLine 2\nLine 3",
			expected: "Line 1\nLine 2\nLine 3",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "text with special characters",
			input:    "You are a code reviewer! Check for: bugs, performance, security.",
			expected: "You are a code reviewer! Check for: bugs, performance, security.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := LoadPrompt(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestLoadPrompt_File(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Test cases
	tests := []struct {
		name        string
		fileContent string
		fileName    string
		useFileURI  bool
		wantError   bool
	}{
		{
			name:        "simple file",
			fileContent: "You are a helpful assistant from file",
			fileName:    "prompt.txt",
			useFileURI:  false,
			wantError:   false,
		},
		{
			name:        "file with file:// prefix",
			fileContent: "You are a helpful assistant with file:// prefix",
			fileName:    "prompt2.txt",
			useFileURI:  true,
			wantError:   false,
		},
		{
			name:        "multiline file",
			fileContent: "Line 1\nLine 2\nLine 3",
			fileName:    "multiline.txt",
			useFileURI:  false,
			wantError:   false,
		},
		{
			name:        "file not found with file:// prefix",
			fileContent: "",
			fileName:    "nonexistent.txt",
			useFileURI:  true,
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join(tmpDir, tt.fileName)

			// Create file if content is provided
			if tt.fileContent != "" {
				err := os.WriteFile(filePath, []byte(tt.fileContent), 0o600)
				if err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			}

			// Prepare input with or without file:// prefix
			input := filePath
			if tt.useFileURI {
				input = "file://" + filePath
			}

			// Test LoadPrompt
			result, err := LoadPrompt(input)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result != tt.fileContent {
				t.Errorf("expected %q, got %q", tt.fileContent, result)
			}
		})
	}
}

func TestLoadPrompt_URL(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		statusCode  int
		wantError   bool
		errorString string
	}{
		{
			name:       "successful fetch",
			content:    "You are a helpful assistant from URL",
			statusCode: http.StatusOK,
			wantError:  false,
		},
		{
			name:        "404 not found",
			content:     "",
			statusCode:  http.StatusNotFound,
			wantError:   true,
			errorString: "status code 404",
		},
		{
			name:        "500 server error",
			content:     "",
			statusCode:  http.StatusInternalServerError,
			wantError:   true,
			errorString: "status code 500",
		},
		{
			name:       "multiline content",
			content:    "Line 1\nLine 2\nLine 3",
			statusCode: http.StatusOK,
			wantError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Check User-Agent header
					if userAgent := r.Header.Get("User-Agent"); userAgent != "LLM-Action/1.0" {
						t.Errorf("expected User-Agent 'LLM-Action/1.0', got %q", userAgent)
					}

					w.WriteHeader(tt.statusCode)
					if tt.statusCode == http.StatusOK {
						if _, err := w.Write([]byte(tt.content)); err != nil {
							t.Errorf("failed to write response: %v", err)
						}
					}
				}),
			)
			defer server.Close()

			// Test LoadPrompt
			result, err := LoadPrompt(server.URL)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.errorString != "" && !strings.Contains(err.Error(), tt.errorString) {
					t.Errorf("expected error to contain %q, got %q", tt.errorString, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result != tt.content {
				t.Errorf("expected %q, got %q", tt.content, result)
			}
		})
	}
}

func TestLoadPrompt_InvalidURL(t *testing.T) {
	// Test with invalid URL
	_, err := LoadPrompt("http://invalid-domain-that-does-not-exist-12345.com/prompt.txt")
	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestIsURL(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"http://example.com", true},
		{"https://example.com", true},
		{"https://example.com/path/to/file.txt", true},
		{"file://path/to/file", false},
		{"/path/to/file", false},
		{"plain text", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isURL(tt.input)
			if result != tt.expected {
				t.Errorf("isURL(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsFilePath(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(tmpFile, []byte("test"), 0o600)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "existing file",
			input:    tmpFile,
			expected: true,
		},
		{
			name:     "file:// with existing file",
			input:    "file://" + tmpFile,
			expected: true,
		},
		{
			name:     "file:// with non-existing file",
			input:    "file:///non/existent/file.txt",
			expected: true, // Returns true because of file:// prefix
		},
		{
			name:     "non-existing file without file://",
			input:    "/non/existent/file.txt",
			expected: false,
		},
		{
			name:     "URL",
			input:    "https://example.com/file.txt",
			expected: false,
		},
		{
			name:     "plain text",
			input:    "You are a helpful assistant",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isFilePath(tt.input)
			if result != tt.expected {
				t.Errorf("isFilePath(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLoadPrompt_WithTemplate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		envVars  map[string]string
		expected string
		wantErr  bool
	}{
		{
			name:     "plain text with template variable",
			input:    "Hello, {{.NAME}}!",
			envVars:  map[string]string{"NAME": "World"},
			expected: "Hello, World!",
			wantErr:  false,
		},
		{
			name:     "template with INPUT_ prefix removal",
			input:    "Model: {{.MODEL}}, Temperature: {{.TEMPERATURE}}",
			envVars:  map[string]string{"INPUT_MODEL": "gpt-4", "INPUT_TEMPERATURE": "0.7"},
			expected: "Model: gpt-4, Temperature: 0.7",
			wantErr:  false,
		},
		{
			name:  "template with GitHub Actions variables",
			input: "Analyzing {{.GITHUB_REPOSITORY}} on branch {{.GITHUB_REF_NAME}}",
			envVars: map[string]string{
				"GITHUB_REPOSITORY": "owner/repo",
				"GITHUB_REF_NAME":   "main",
			},
			expected: "Analyzing owner/repo on branch main",
			wantErr:  false,
		},
		{
			name:     "template with conditional",
			input:    "{{if .DEBUG}}Debug enabled{{else}}Debug disabled{{end}}",
			envVars:  map[string]string{"DEBUG": "true"},
			expected: "Debug enabled",
			wantErr:  false,
		},
		{
			name:     "invalid template syntax",
			input:    "{{.VAR",
			envVars:  map[string]string{},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}
			defer func() {
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			result, err := LoadPrompt(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadPrompt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("LoadPrompt() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestLoadPrompt_FileWithTemplate(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Create a file with template content
	templateContent := "Repository: {{.GITHUB_REPOSITORY}}\nModel: {{.MODEL}}"
	filePath := filepath.Join(tmpDir, "template.txt")
	err := os.WriteFile(filePath, []byte(templateContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Set up environment variables
	os.Setenv("GITHUB_REPOSITORY", "appleboy/LLM-action")
	os.Setenv("INPUT_MODEL", "gpt-4o")
	defer func() {
		os.Unsetenv("GITHUB_REPOSITORY")
		os.Unsetenv("INPUT_MODEL")
	}()

	// Test LoadPrompt with file path
	result, err := LoadPrompt(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Repository: appleboy/LLM-action\nModel: gpt-4o"
	if result != expected {
		t.Errorf("LoadPrompt() = %q, want %q", result, expected)
	}
}

func TestLoadPrompt_URLWithTemplate(t *testing.T) {
	// Create test server that returns template content
	templateContent := "URL Test: {{.TEST_VAR}}"
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(templateContent)); err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		}),
	)
	defer server.Close()

	// Set up environment variable
	os.Setenv("TEST_VAR", "success")
	defer os.Unsetenv("TEST_VAR")

	// Test LoadPrompt with URL
	result, err := LoadPrompt(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "URL Test: success"
	if result != expected {
		t.Errorf("LoadPrompt() = %q, want %q", result, expected)
	}
}

// Tests for LoadContent function (without template rendering)

func TestLoadContent_PlainText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple text",
			input:    "-----BEGIN CERTIFICATE-----\nMIIB...\n-----END CERTIFICATE-----",
			expected: "-----BEGIN CERTIFICATE-----\nMIIB...\n-----END CERTIFICATE-----",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "text with template syntax should not be rendered",
			input:    "{{.SHOULD_NOT_RENDER}}",
			expected: "{{.SHOULD_NOT_RENDER}}",
		},
	}

	// Set up environment variable that should NOT be rendered
	os.Setenv("SHOULD_NOT_RENDER", "rendered")
	defer os.Unsetenv("SHOULD_NOT_RENDER")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := LoadContent(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestLoadContent_File(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Create test file with certificate-like content
	certContent := "-----BEGIN CERTIFICATE-----\nMIIDxTCCAq2gAwIBAgIQAqx...\n-----END CERTIFICATE-----"
	filePath := filepath.Join(tmpDir, "ca-cert.pem")
	err := os.WriteFile(filePath, []byte(certContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	tests := []struct {
		name      string
		input     string
		expected  string
		wantError bool
	}{
		{
			name:      "file path",
			input:     filePath,
			expected:  certContent,
			wantError: false,
		},
		{
			name:      "file:// prefix",
			input:     "file://" + filePath,
			expected:  certContent,
			wantError: false,
		},
		{
			name:      "non-existent file with file:// prefix",
			input:     "file:///non/existent/file.pem",
			expected:  "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := LoadContent(tt.input)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestLoadContent_URL(t *testing.T) {
	certContent := "-----BEGIN CERTIFICATE-----\nMIIDxTCCAq2gAwIBAgIQAqx...\n-----END CERTIFICATE-----"

	tests := []struct {
		name       string
		content    string
		statusCode int
		wantError  bool
	}{
		{
			name:       "successful fetch",
			content:    certContent,
			statusCode: http.StatusOK,
			wantError:  false,
		},
		{
			name:       "404 not found",
			content:    "",
			statusCode: http.StatusNotFound,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.statusCode)
					if tt.statusCode == http.StatusOK {
						if _, err := w.Write([]byte(tt.content)); err != nil {
							t.Errorf("failed to write response: %v", err)
						}
					}
				}),
			)
			defer server.Close()

			result, err := LoadContent(server.URL)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result != tt.content {
				t.Errorf("expected %q, got %q", tt.content, result)
			}
		})
	}
}

func TestLoadContent_NoTemplateRendering(t *testing.T) {
	// Create a temporary file with template syntax
	tmpDir := t.TempDir()
	templateContent := "Content with {{.VAR}} that should not be rendered"
	filePath := filepath.Join(tmpDir, "template.txt")
	err := os.WriteFile(filePath, []byte(templateContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Set environment variable
	os.Setenv("VAR", "should-not-appear")
	defer os.Unsetenv("VAR")

	// LoadContent should NOT render the template
	result, err := LoadContent(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The template syntax should remain as-is
	if result != templateContent {
		t.Errorf(
			"LoadContent() should not render templates, got %q, want %q",
			result,
			templateContent,
		)
	}
}
