# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

LLM Action is a GitHub Action that enables interaction with OpenAI-compatible LLM services. It supports any OpenAI-compatible API endpoint including OpenAI, Azure OpenAI, Ollama, LocalAI, LM Studio, vLLM, and other self-hosted services.

## Development Commands

### Testing

```bash
# Run all tests with race detection and coverage
go test -race -cover -coverprofile=coverage.out ./...

# Run a specific test
go test -v -run TestName ./...
```

### Linting

```bash
# Run golangci-lint (requires golangci-lint v2.6)
golangci-lint run --verbose

# Check Dockerfile
hadolint Dockerfile

# Fix the golang format
golangci-lint fmt
```

### Building

```bash
# Build the binary locally
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o llm-action .

# Build Docker image
docker build -t llm-action .
```

### Running Locally

The action reads configuration from environment variables with `INPUT_` prefix:

```bash
export INPUT_API_KEY="your-api-key"
export INPUT_INPUT_PROMPT="your prompt"
export INPUT_MODEL="gpt-4o"  # optional
export INPUT_BASE_URL="https://api.openai.com/v1"  # optional
export INPUT_DEBUG="true"  # optional
go run .
```

## Architecture

### Core Components

**main.go** - Entry point and orchestration

- `run()` function orchestrates the entire flow: config loading → client creation → message building → API call → output handling
- `maskAPIKey()` securely masks API keys in debug output (shows first/last 4 chars)
- Uses `github.com/appleboy/com/gh.SetOutput()` to set GitHub Actions outputs

**config.go** - Configuration management

- `LoadConfig()` reads all inputs from environment variables (GitHub Actions sets these with `INPUT_` prefix)
- Required: `api_key`, `input_prompt`
- Optional with defaults: `base_url`, `model`, `temperature`, `max_tokens`, `skip_ssl_verify`, `system_prompt`, `debug`
- Each optional parameter has dedicated parse methods (`parseTemperature`, `parseMaxTokens`, `parseSkipSSL`, `parseDebug`)

**client.go** - OpenAI client initialization

- `NewClient()` creates configured OpenAI client from `github.com/sashabaranov/go-openai`
- `createInsecureHTTPClient()` creates HTTP client with SSL verification disabled (marked with `#nosec G402` for gosec)
- SSL skip is intentionally configurable for local/self-hosted LLM services

**message.go** - Message construction

- `BuildMessages()` constructs OpenAI chat completion message array
- System prompt is prepended if provided, followed by user input prompt
- Returns slice of `openai.ChatCompletionMessage`

### Data Flow

1. Environment variables (`INPUT_*`) → `LoadConfig()` → `Config` struct
2. `Config` → `NewClient()` → OpenAI client with custom base URL and SSL settings
3. `Config` → `BuildMessages()` → OpenAI message format
4. Client + Messages → `CreateChatCompletion()` → API call to LLM service
5. API response → Extract content → `gh.SetOutput()` → GitHub Actions output

### GitHub Actions Integration

- Defined in `action.yml` with inputs/outputs specification
- Runs using Docker (multi-stage build in `Dockerfile`)
- Docker image uses non-root user (appuser:1000) for security
- Go binary is statically compiled (`CGO_ENABLED=0`) for Alpine Linux base image

## Code Style & Standards

- Go version: 1.25
- Linter configuration in `.golangci.yml` includes: gosec, govet, staticcheck, errcheck, and formatting tools (gofmt, gofumpt, goimports, golines)
- Security scanning via gosec is enabled; intentional security exceptions are marked with `#nosec` comments
- All tests should include race detection (`-race` flag)

## Testing Philosophy

- Table-driven tests for configuration parsing (see `config_test.go`)
- Mock/test implementations for message building (see `message_test.go`)
- Test coverage is uploaded to Codecov via CI

## Deployment

- GitHub Actions workflows in `.github/workflows/`:
  - `testing.yml` - Runs tests and linting on push/PR
  - `trivy.yml` - Security scanning
  - `docker.yml` - Docker image building
  - `goreleaser.yml` - Release automation
  - `codeql.yml` - Code security analysis

## Key Dependencies

- `github.com/sashabaranov/go-openai` - OpenAI API client library
- `github.com/appleboy/com` - GitHub Actions helper utilities for output handling
- `github.com/yassinebenaid/godump` - Pretty printing for debug mode
