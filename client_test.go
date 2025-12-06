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
		expectNil     bool
		expectError   bool
	}{
		{
			name:          "No custom config returns nil",
			caCert:        "",
			skipSSLVerify: false,
			expectNil:     true,
			expectError:   false,
		},
		{
			name:          "Skip SSL only",
			caCert:        "",
			skipSSLVerify: true,
			expectNil:     false,
			expectError:   false,
		},
		{
			name:          "Invalid CA certificate",
			caCert:        "invalid-cert",
			skipSSLVerify: false,
			expectNil:     false,
			expectError:   true,
		},
		{
			name:          "Invalid CA certificate with skip SSL",
			caCert:        "invalid-cert",
			skipSSLVerify: true,
			expectNil:     false,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := createHTTPClient(tt.caCert, tt.skipSSLVerify)

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

			if tt.expectNil {
				if client != nil {
					t.Error("expected nil client, got non-nil")
				}
				return
			}

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
