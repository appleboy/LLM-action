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
