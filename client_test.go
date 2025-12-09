package main

import (
	"testing"
)

// Sample valid CA certificate for testing (self-signed)
const testCACert = `-----BEGIN CERTIFICATE-----
MIIBkTCB+wIJAKHBfpFoSdLGMA0GCSqGSIb3DQEBCwUAMBExDzANBgNVBAMMBnRl
c3RjYTAeFw0yMzAxMDEwMDAwMDBaFw0yNDAxMDEwMDAwMDBaMBExDzANBgNVBAMM
BnRlc3RjYTBcMA0GCSqGSIb3DQEBAQUAA0sAMEgCQQC5Q5mJKL8nKL8nKL8nKL8n
KL8nKL8nKL8nKL8nKL8nKL8nKL8nKL8nKL8nKL8nKL8nKL8nKL8nKL8nKL8nKL8n
AgMBAAEwDQYJKoZIhvcNAQELBQADQQBRvYFpCgK/G8K/G8K/G8K/G8K/G8K/G8K/
G8K/G8K/G8K/G8K/G8K/G8K/G8K/G8K/G8K/G8K/G8K/G8K/G8K/G8K/
-----END CERTIFICATE-----`

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name: "Default OpenAI config",
			config: &Config{
				BaseURL:       "https://api.openai.com/v1",
				APIKey:        "test-key",
				SkipSSLVerify: false,
			},
			expectError: false,
		},
		{
			name: "Custom base URL",
			config: &Config{
				BaseURL:       "http://localhost:8080/v1",
				APIKey:        "test-key",
				SkipSSLVerify: false,
			},
			expectError: false,
		},
		{
			name: "Skip SSL verification",
			config: &Config{
				BaseURL:       "https://api.openai.com/v1",
				APIKey:        "test-key",
				SkipSSLVerify: true,
			},
			expectError: false,
		},
		{
			name: "With invalid CA certificate",
			config: &Config{
				BaseURL:       "https://api.openai.com/v1",
				APIKey:        "test-key",
				CACert:        "invalid-cert",
				SkipSSLVerify: false,
			},
			expectError: true,
		},
		{
			name: "Skip SSL with CA certificate",
			config: &Config{
				BaseURL:       "https://api.openai.com/v1",
				APIKey:        "test-key",
				CACert:        testCACert,
				SkipSSLVerify: true,
			},
			expectError: true, // Invalid test cert will fail
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if client == nil {
				t.Error("expected client to be created, got nil")
			}
		})
	}
}

func TestCreateHTTPClient(t *testing.T) {
	tests := []struct {
		name          string
		caCert        string
		skipSSLVerify bool
		customHeaders map[string]string
		expectError   bool
	}{
		{
			name:          "Default config with default headers",
			caCert:        "",
			skipSSLVerify: false,
			customHeaders: nil,
			expectError:   false,
		},
		{
			name:          "Skip SSL only",
			caCert:        "",
			skipSSLVerify: true,
			customHeaders: nil,
			expectError:   false,
		},
		{
			name:          "Invalid CA certificate",
			caCert:        "invalid-cert",
			skipSSLVerify: false,
			customHeaders: nil,
			expectError:   true,
		},
		{
			name:          "Invalid CA certificate with skip SSL",
			caCert:        "invalid-cert",
			skipSSLVerify: true,
			customHeaders: nil,
			expectError:   true,
		},
		{
			name:          "Custom headers only",
			caCert:        "",
			skipSSLVerify: false,
			customHeaders: map[string]string{"X-Custom": "value"},
			expectError:   false,
		},
		{
			name:          "Custom headers with skip SSL",
			caCert:        "",
			skipSSLVerify: true,
			customHeaders: map[string]string{"X-Request-ID": "test123"},
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := createHTTPClient(tt.caCert, tt.skipSSLVerify, tt.customHeaders)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Client should always be non-nil (default headers are always added)
			if client == nil {
				t.Error("expected non-nil client, got nil")
				return
			}

			if client.Transport == nil {
				t.Error("expected transport to be set, got nil")
			}
		})
	}
}

func TestGetDefaultHeaders(t *testing.T) {
	headers := getDefaultHeaders()

	if headers["User-Agent"] == "" {
		t.Error("expected User-Agent header to be set")
	}
	if headers["X-Action-Name"] != ActionName {
		t.Errorf("expected X-Action-Name to be %s, got %s", ActionName, headers["X-Action-Name"])
	}
	if headers["X-Action-Version"] == "" {
		t.Error("expected X-Action-Version header to be set")
	}
}

func TestMergeHeaders(t *testing.T) {
	customHeaders := map[string]string{
		"X-Custom":   "value",
		"User-Agent": "custom-agent", // Override default
	}

	merged := mergeHeaders(customHeaders)

	// Custom headers should be present
	if merged["X-Custom"] != "value" {
		t.Errorf("expected X-Custom to be 'value', got '%s'", merged["X-Custom"])
	}

	// Custom User-Agent should override default
	if merged["User-Agent"] != "custom-agent" {
		t.Errorf("expected User-Agent to be 'custom-agent', got '%s'", merged["User-Agent"])
	}

	// Default headers should still be present
	if merged["X-Action-Name"] != ActionName {
		t.Errorf("expected X-Action-Name to be %s, got %s", ActionName, merged["X-Action-Name"])
	}
}
