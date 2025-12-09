package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	openai "github.com/sashabaranov/go-openai"
)

// headerTransport wraps an http.RoundTripper to add custom headers to requests
type headerTransport struct {
	base    http.RoundTripper
	headers map[string]string
}

// RoundTrip implements http.RoundTripper interface
func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	reqClone := req.Clone(req.Context())
	for key, value := range t.headers {
		reqClone.Header.Set(key, value)
	}
	return t.base.RoundTrip(reqClone)
}

// NewClient creates a new OpenAI client with the given configuration
func NewClient(config *Config) (*openai.Client, error) {
	clientConfig := openai.DefaultConfig(config.APIKey)
	clientConfig.BaseURL = config.BaseURL

	// Handle custom CA certificate, SSL verification, and headers
	httpClient, err := createHTTPClient(config.CACert, config.SkipSSLVerify, config.Headers)
	if err != nil {
		return nil, err
	}
	clientConfig.HTTPClient = httpClient

	return openai.NewClientWithConfig(clientConfig), nil
}

// getDefaultHeaders returns the default headers for all API requests
func getDefaultHeaders() map[string]string {
	return map[string]string{
		"User-Agent":       GetUserAgent(),
		"X-Action-Name":    ActionName,
		"X-Action-Version": Version,
	}
}

// mergeHeaders merges custom headers with default headers.
// Custom headers take precedence over default headers.
func mergeHeaders(customHeaders map[string]string) map[string]string {
	merged := getDefaultHeaders()
	for key, value := range customHeaders {
		merged[key] = value
	}
	return merged
}

// createHTTPClient creates an HTTP client with optional custom CA certificate,
// SSL verification settings, and headers (including default action headers).
func createHTTPClient(
	caCert string,
	skipSSLVerify bool,
	customHeaders map[string]string,
) (*http.Client, error) {
	baseTransport := http.DefaultTransport

	// Only create custom TLS transport if needed
	if caCert != "" || skipSSLVerify {
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

		baseTransport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	}

	// Always wrap transport with headers (default + custom)
	allHeaders := mergeHeaders(customHeaders)
	finalTransport := &headerTransport{
		base:    baseTransport,
		headers: allHeaders,
	}

	return &http.Client{
		Transport: finalTransport,
	}, nil
}
