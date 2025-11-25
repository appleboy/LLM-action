# LLM Action

[English](README.md) | [繁體中文](README.zh-TW.md) | [簡體中文](README.zh-CN.md)

[![Lint and Testing](https://github.com/appleboy/LLM-action/actions/workflows/testing.yml/badge.svg)](https://github.com/appleboy/LLM-action/actions/workflows/testing.yml)
[![Trivy Security Scan](https://github.com/appleboy/LLM-action/actions/workflows/trivy.yml/badge.svg)](https://github.com/appleboy/LLM-action/actions/workflows/trivy.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/LLM-action)](https://goreportcard.com/report/github.com/appleboy/LLM-action)

A GitHub Action to interact with OpenAI Compatible LLM services. This action allows you to connect to any OpenAI-compatible API endpoint (including local or self-hosted services) and get responses that can be used in your workflow.

## Features

- 🔌 Connect to any OpenAI-compatible API endpoint
- 🔐 Support for custom API keys
- 🔧 Configurable base URL for self-hosted services
- 🚫 Optional SSL certificate verification skip
- 🎯 System prompt support for context setting
- 📝 Output response available for subsequent actions
- 🎛️ Configurable temperature and max tokens
- 🐛 Debug mode with secure API key masking

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `base_url` | Base URL for OpenAI Compatible API endpoint | No | `https://api.openai.com/v1` |
| `api_key` | API Key for authentication | Yes | - |
| `model` | Model name to use | No | `gpt-4o` |
| `skip_ssl_verify` | Skip SSL certificate verification | No | `false` |
| `system_prompt` | System prompt to set the context. Supports plain text, file path, or URL | No | `''` |
| `input_prompt` | User input prompt for the LLM | Yes | - |
| `temperature` | Temperature for response randomness (0.0-2.0) | No | `0.7` |
| `max_tokens` | Maximum tokens in the response | No | `1000` |
| `debug` | Enable debug mode to print all parameters (API key will be masked) | No | `false` |

## Outputs

| Output | Description |
|--------|-------------|
| `response` | The response from the LLM |

## Usage Examples

### Basic Example

```yaml
name: LLM Workflow
on: [push]

jobs:
  llm-task:
    runs-on: ubuntu-latest
    steps:
      - name: Call LLM
        id: llm
        uses: appleboy/LLM-action@v1
        with:
          api_key: ${{ secrets.OPENAI_API_KEY }}
          input_prompt: 'What is GitHub Actions?'

      - name: Use LLM Response
        run: |
          echo "LLM Response:"
          echo "${{ steps.llm.outputs.response }}"
```

### With System Prompt

```yaml
- name: Code Review with LLM
  id: review
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    model: 'gpt-4'
    system_prompt: 'You are a code reviewer. Provide constructive feedback on code quality, best practices, and potential issues.'
    input_prompt: |
      Review this code:
      ```python
      def add(a, b):
          return a + b
      ```
    temperature: '0.3'
    max_tokens: '2000'

- name: Post Review Comment
  run: |
    echo "${{ steps.review.outputs.response }}"
```

### With Multiline System Prompt

```yaml
- name: Advanced Code Review
  id: review
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    model: 'gpt-4'
    system_prompt: |
      You are an expert code reviewer with deep knowledge of software engineering best practices.

      Your responsibilities:
      - Identify potential bugs and security vulnerabilities
      - Suggest improvements for code quality and maintainability
      - Check for adherence to coding standards
      - Evaluate performance implications

      Provide constructive, actionable feedback in a professional tone.
    input_prompt: |
      Review the following pull request changes:
      ${{ github.event.pull_request.body }}
    temperature: '0.3'
    max_tokens: '2000'
```

### System Prompt from File

Instead of embedding long prompts in YAML, you can load them from a file:

```yaml
- name: Code Review with Prompt File
  id: review
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    model: 'gpt-4'
    system_prompt: '.github/prompts/code-review.txt'
    input_prompt: |
      Review this code:
      ```python
      def calculate(x, y):
          return x / y
      ```
```

Or using `file://` prefix:

```yaml
- name: Code Review with File URI
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    system_prompt: 'file://.github/prompts/code-review.txt'
    input_prompt: 'Review the main.go file'
```

### System Prompt from URL

Load prompts from a remote URL:

```yaml
- name: Code Review with Remote Prompt
  id: review
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    model: 'gpt-4'
    system_prompt: 'https://raw.githubusercontent.com/your-org/prompts/main/code-review.txt'
    input_prompt: |
      Review this pull request:
      ${{ github.event.pull_request.body }}
```

### Self-Hosted / Local LLM

```yaml
- name: Call Local LLM
  id: local_llm
  uses: appleboy/LLM-action@v1
  with:
    base_url: 'http://localhost:8080/v1'
    api_key: 'your-local-api-key'
    model: 'llama2'
    skip_ssl_verify: 'true'
    input_prompt: 'Explain quantum computing in simple terms'
```

### Using with Ollama

```yaml
- name: Call Ollama
  id: ollama
  uses: appleboy/LLM-action@v1
  with:
    base_url: 'http://localhost:11434/v1'
    api_key: 'ollama'
    model: 'llama3'
    system_prompt: 'You are a helpful assistant'
    input_prompt: 'Write a haiku about programming'
```

### Chain Multiple LLM Calls

```yaml
- name: Generate Story
  id: generate
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    input_prompt: 'Write a short story about a robot'
    max_tokens: '500'

- name: Translate Story
  id: translate
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    system_prompt: 'You are a translator'
    input_prompt: |
      Translate the following text to Spanish:
      ${{ steps.generate.outputs.response }}

- name: Display Results
  run: |
    echo "Original Story:"
    echo "${{ steps.generate.outputs.response }}"
    echo ""
    echo "Translated Story:"
    echo "${{ steps.translate.outputs.response }}"
```

### Debug Mode

Enable debug mode to troubleshoot issues and inspect all parameters:

```yaml
- name: Call LLM with Debug
  id: llm_debug
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    model: 'gpt-4'
    system_prompt: 'You are a helpful assistant'
    input_prompt: 'Explain how GitHub Actions work'
    temperature: '0.8'
    max_tokens: '1500'
    debug: true  # Enable debug mode
```

**Debug Output Example:**

```txt
=== Debug Mode: All Parameters ===
main.Config{
    BaseURL: "https://api.openai.com/v1",
    APIKey: "sk-ab****xyz9",  // Masked for security
    Model: "gpt-4",
    SkipSSLVerify: false,
    SystemPrompt: "You are a helpful assistant",
    InputPrompt: "Explain how GitHub Actions work",
    Temperature: 0.8,
    MaxTokens: 1500,
    Debug: true
}
===================================
=== Debug Mode: Messages ===
[... message details ...]
============================
```

**Security Note:** When debug mode is enabled, the API key is automatically masked (only showing first 4 and last 4 characters) to prevent accidental exposure in logs.

## Supported Services

This action works with any OpenAI-compatible API, including:

- **OpenAI** - `https://api.openai.com/v1`
- **Azure OpenAI** - `https://{your-resource}.openai.azure.com/openai/deployments/{deployment-id}`
- **Ollama** - `http://localhost:11434/v1`
- **LocalAI** - `http://localhost:8080/v1`
- **LM Studio** - `http://localhost:1234/v1`
- **Jan** - `http://localhost:1337/v1`
- **vLLM** - Your vLLM server endpoint
- **Text Generation WebUI** - Your WebUI endpoint
- Any other OpenAI-compatible service

## Security Considerations

- Always use GitHub Secrets for API keys: `${{ secrets.YOUR_API_KEY }}`
- Only use `skip_ssl_verify: 'true'` for trusted local/internal services
- Be careful with sensitive data in prompts, as they will be sent to the LLM service

## License

MIT License - see LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
