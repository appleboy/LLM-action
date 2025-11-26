package main

import (
	"os"
	"strings"
	"testing"
)

const testModelGPT4 = "gpt-4"

func TestRenderTemplate(t *testing.T) {
	tests := []struct {
		name     string
		template string
		envVars  map[string]string
		want     string
		wantErr  bool
	}{
		{
			name:     "plain text without template",
			template: "Hello, World!",
			envVars:  map[string]string{},
			want:     "Hello, World!",
			wantErr:  false,
		},
		{
			name:     "simple variable substitution",
			template: "Hello, {{.NAME}}!",
			envVars:  map[string]string{"NAME": "Alice"},
			want:     "Hello, Alice!",
			wantErr:  false,
		},
		{
			name:     "INPUT_ prefix removal",
			template: "Model: {{.MODEL}}, API Key: {{.API_KEY}}",
			envVars: map[string]string{
				"INPUT_MODEL":   testModelGPT4,
				"INPUT_API_KEY": "sk-xxx",
			},
			want:    "Model: gpt-4, API Key: sk-xxx",
			wantErr: false,
		},
		{
			name:     "INPUT_ prefix with both forms",
			template: "With prefix: {{.INPUT_MODEL}}, Without: {{.MODEL}}",
			envVars:  map[string]string{"INPUT_MODEL": testModelGPT4},
			want:     "With prefix: gpt-4, Without: gpt-4",
			wantErr:  false,
		},
		{
			name:     "multiple variables",
			template: "Repo: {{.GITHUB_REPOSITORY}}, Branch: {{.GITHUB_REF}}, Model: {{.MODEL}}",
			envVars: map[string]string{
				"GITHUB_REPOSITORY": "owner/repo",
				"GITHUB_REF":        "refs/heads/main",
				"INPUT_MODEL":       "gpt-4o",
			},
			want:    "Repo: owner/repo, Branch: refs/heads/main, Model: gpt-4o",
			wantErr: false,
		},
		{
			name:     "template with conditionals",
			template: "{{if .DEBUG}}Debug mode enabled{{else}}Debug mode disabled{{end}}",
			envVars:  map[string]string{"DEBUG": "true"},
			want:     "Debug mode enabled",
			wantErr:  false,
		},
		{
			name:     "template with missing variable",
			template: "Hello, {{.MISSING_VAR}}!",
			envVars:  map[string]string{},
			want:     "Hello, <no value>!",
			wantErr:  false,
		},
		{
			name:     "invalid template syntax",
			template: "Hello, {{.NAME}!",
			envVars:  map[string]string{"NAME": "Alice"},
			want:     "",
			wantErr:  true,
		},
		{
			name:     "template with multiline",
			template: "Line 1: {{.VAR1}}\nLine 2: {{.VAR2}}",
			envVars: map[string]string{
				"VAR1": "value1",
				"VAR2": "value2",
			},
			want:    "Line 1: value1\nLine 2: value2",
			wantErr: false,
		},
		{
			name: "template with GitHub Actions variables",
			template: `Please analyze the repository: {{.GITHUB_REPOSITORY}}
Using model: {{.MODEL}}
Current branch: {{.GITHUB_REF_NAME}}`,
			envVars: map[string]string{
				"GITHUB_REPOSITORY": "appleboy/LLM-action",
				"INPUT_MODEL":       "gpt-4o",
				"GITHUB_REF_NAME":   "main",
			},
			want: `Please analyze the repository: appleboy/LLM-action
Using model: gpt-4o
Current branch: main`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}
			defer func() {
				// Clean up environment variables
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			got, err := RenderTemplate(tt.template)
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RenderTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildTemplateData(t *testing.T) {
	// Set up test environment variables
	testEnvVars := map[string]string{
		"INPUT_MODEL":       testModelGPT4,
		"INPUT_API_KEY":     "sk-test",
		"GITHUB_REPOSITORY": "owner/repo",
		"PATH":              "/usr/bin",
	}

	for key, value := range testEnvVars {
		os.Setenv(key, value)
	}
	defer func() {
		for key := range testEnvVars {
			os.Unsetenv(key)
		}
	}()

	data := buildTemplateData()

	// Check that INPUT_ prefixed variables are available both ways
	if data["INPUT_MODEL"] != testModelGPT4 {
		t.Errorf("Expected INPUT_MODEL to be %s, got %s", testModelGPT4, data["INPUT_MODEL"])
	}
	if data["MODEL"] != testModelGPT4 {
		t.Errorf("Expected MODEL (without prefix) to be %s, got %s", testModelGPT4, data["MODEL"])
	}

	if data["INPUT_API_KEY"] != "sk-test" {
		t.Errorf("Expected INPUT_API_KEY to be sk-test, got %s", data["INPUT_API_KEY"])
	}
	if data["API_KEY"] != "sk-test" {
		t.Errorf("Expected API_KEY (without prefix) to be sk-test, got %s", data["API_KEY"])
	}

	// Check that non-INPUT_ variables are only available with original key
	if data["GITHUB_REPOSITORY"] != "owner/repo" {
		t.Errorf("Expected GITHUB_REPOSITORY to be owner/repo, got %s", data["GITHUB_REPOSITORY"])
	}

	if data["PATH"] != "/usr/bin" {
		t.Errorf("Expected PATH to be /usr/bin, got %s", data["PATH"])
	}
}

func TestRenderTemplateWithRealEnv(t *testing.T) {
	// Set up environment
	os.Setenv("INPUT_TEST_VAR", "test_value")
	os.Setenv("NORMAL_VAR", "normal_value")
	defer func() {
		os.Unsetenv("INPUT_TEST_VAR")
		os.Unsetenv("NORMAL_VAR")
	}()

	tests := []struct {
		name     string
		template string
		want     string
	}{
		{
			name:     "access INPUT_ variable without prefix",
			template: "Value: {{.TEST_VAR}}",
			want:     "Value: test_value",
		},
		{
			name:     "access INPUT_ variable with prefix",
			template: "Value: {{.INPUT_TEST_VAR}}",
			want:     "Value: test_value",
		},
		{
			name:     "access normal variable",
			template: "Value: {{.NORMAL_VAR}}",
			want:     "Value: normal_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RenderTemplate(tt.template)
			if err != nil {
				t.Errorf("RenderTemplate() unexpected error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("RenderTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenderTemplateErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		template string
	}{
		{
			name:     "unclosed action",
			template: "{{.VAR",
		},
		{
			name:     "invalid action",
			template: "{{.VAR)}}",
		},
		{
			name:     "unclosed conditional",
			template: "{{if .VAR}}text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := RenderTemplate(tt.template)
			if err == nil {
				t.Errorf("RenderTemplate() expected error for invalid template, got nil")
			}
			if !strings.Contains(err.Error(), "failed to") {
				t.Errorf("RenderTemplate() error message should contain 'failed to', got: %v", err)
			}
		})
	}
}
