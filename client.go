package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	openai "github.com/sashabaranov/go-openai"
)

// NewClient creates a new OpenAI client with the given configuration
func NewClient(config *Config) (*openai.Client, error) {
	clientConfig := openai.DefaultConfig(config.APIKey)
	clientConfig.BaseURL = config.BaseURL

	// Handle custom CA certificate and SSL verification
	httpClient, err := createHTTPClient(config.CACert, config.SkipSSLVerify)
	if err != nil {
		return nil, err
	}
	if httpClient != nil {
		clientConfig.HTTPClient = httpClient
	}

	return openai.NewClientWithConfig(clientConfig), nil
}

// createHTTPClient creates an HTTP client with optional custom CA certificate and SSL verification settings
// Returns nil if no custom configuration is needed
func createHTTPClient(caCert string, skipSSLVerify bool) (*http.Client, error) {
	// If no custom CA cert and SSL verification is enabled, use default client
	if caCert == "" && !skipSSLVerify {
		return nil, nil
	}

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Handle custom CA certificate
	if caCert != "" {
		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM([]byte(caCert)); !ok {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}
		tlsConfig.RootCAs = caCertPool
	}

	// Handle SSL verification skip
	// #nosec G402 - This is intentionally configurable by the user for local/self-hosted LLM services
	if skipSSLVerify {
		tlsConfig.InsecureSkipVerify = true
	}

	customTransport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return &http.Client{
		Transport: customTransport,
	}, nil
}
