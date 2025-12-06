# LLM Action

[English](README.md) | [繁體中文](README.zh-TW.md) | [簡體中文](README.zh-CN.md)

[![Lint and Testing](https://github.com/appleboy/LLM-action/actions/workflows/testing.yml/badge.svg)](https://github.com/appleboy/LLM-action/actions/workflows/testing.yml)
[![Trivy Security Scan](https://github.com/appleboy/LLM-action/actions/workflows/trivy.yml/badge.svg)](https://github.com/appleboy/LLM-action/actions/workflows/trivy.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/LLM-action)](https://goreportcard.com/report/github.com/appleboy/LLM-action)

一个用于与 OpenAI 兼容 LLM 服务交互的 GitHub Action。此 Action 允许您连接到任何 OpenAI 兼容的 API 端点（包括本地或自托管服务），并获取可用于工作流程的响应。

## 功能特色

- 🔌 连接任何 OpenAI 兼容的 API 端点
- 🔐 支持自定义 API 密钥
- 🔧 可配置的基础 URL，适用于自托管服务
- 🚫 可选跳过 SSL 证书验证
- 🔒 支持自定义 CA 证书，适用于自签名证书
- 🎯 支持系统提示词以设定上下文
- 📝 输出响应可用于后续 Actions
- 🎛️ 可配置的温度和最大令牌数
- 🐛 调试模式，并安全地屏蔽 API 密钥
- 🎨 支持 Go 模板语法，可动态插入环境变量
- 🛠️ 通过函数调用支持结构化输出（tool schema 支持）

## 输入参数

| 输入              | 说明                                                                            | 必填 | 默认值                      |
| ----------------- | ------------------------------------------------------------------------------- | ---- | --------------------------- |
| `base_url`        | OpenAI 兼容 API 端点的基础 URL                                                  | 否   | `https://api.openai.com/v1` |
| `api_key`         | 用于验证的 API 密钥                                                             | 是   | -                           |
| `model`           | 要使用的模型名称                                                                | 否   | `gpt-4o`                    |
| `skip_ssl_verify` | 跳过 SSL 证书验证                                                               | 否   | `false`                     |
| `ca_cert`         | 自定义 CA 证书。支持证书内容、文件路径或 URL                                    | 否   | `''`                        |
| `system_prompt`   | 设定上下文的系统提示词。支持纯文本、文件路径或 URL。支持 Go 模板语法与环境变量  | 否   | `''`                        |
| `input_prompt`    | 用户输入给 LLM 的提示词。支持纯文本、文件路径或 URL。支持 Go 模板语法与环境变量 | 是   | -                           |
| `tool_schema`     | 用于结构化输出的 JSON schema（函数调用）。支持纯文本、文件路径或 URL。支持 Go 模板语法 | 否   | `''`                        |
| `temperature`     | 响应随机性的温度值（0.0-2.0）                                                   | 否   | `0.7`                       |
| `max_tokens`      | 响应中的最大令牌数                                                              | 否   | `1000`                      |
| `debug`           | 启用调试模式以显示所有参数（API 密钥将被屏蔽）                                  | 否   | `false`                     |

## 输出参数

| 输出       | 说明                                                                     |
| ---------- | ------------------------------------------------------------------------ |
| `response` | 来自 LLM 的响应（未使用 tool_schema 时）                                 |
| `<field>`  | 使用 tool_schema 时，函数参数 JSON 中的每个字段都会成为独立的输出        |

**注意：** 当使用 `tool_schema` 时，输出会根据 schema 的属性动态生成。例如，如果您的 schema 定义了 `city` 和 `country` 字段，输出将会是 `steps.<id>.outputs.city` 和 `steps.<id>.outputs.country`。

## 使用范例

### 基本范例

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
          input_prompt: "What is GitHub Actions?"

      - name: Use LLM Response
        run: |
          echo "LLM Response:"
          echo "${{ steps.llm.outputs.response }}"
```

### 使用系统提示词

````yaml
- name: Code Review with LLM
  id: review
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    model: "gpt-4"
    system_prompt: "你是一位代码审查员。请提供有关代码质量、最佳实践和潜在问题的建设性意见。"
    input_prompt: |
      请审查此代码：
      ```python
      def add(a, b):
          return a + b
      ```
    temperature: "0.3"
    max_tokens: "2000"

- name: Post Review Comment
  run: |
    echo "${{ steps.review.outputs.response }}"
````

### 使用多行系统提示词

```yaml
- name: Advanced Code Review
  id: review
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    model: "gpt-4"
    system_prompt: |
      你是一位拥有深厚软件工程最佳实践知识的专业代码审查员。

      你的职责：
      - 识别潜在的错误和安全漏洞
      - 建议改善代码质量和可维护性的方法
      - 检查是否遵守代码标准
      - 评估性能影响

      请以专业的语气提供建设性、可行的意见。
    input_prompt: |
      审查以下 Pull Request 变更：
      ${{ github.event.pull_request.body }}
    temperature: "0.3"
    max_tokens: "2000"
```

### 从文件加载系统提示词

无需在 YAML 中嵌入冗长的提示词，可以从文件加载：

````yaml
- name: Code Review with Prompt File
  id: review
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    model: "gpt-4"
    system_prompt: ".github/prompts/code-review.txt"
    input_prompt: |
      审查此代码：
      ```python
      def calculate(x, y):
          return x / y
      ```
````

或使用 `file://` 前缀：

```yaml
- name: Code Review with File URI
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    system_prompt: "file://.github/prompts/code-review.txt"
    input_prompt: "审查 main.go 文件"
```

### 从 URL 加载系统提示词

从远程 URL 加载提示词：

```yaml
- name: Code Review with Remote Prompt
  id: review
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    model: "gpt-4"
    system_prompt: "https://raw.githubusercontent.com/your-org/prompts/main/code-review.txt"
    input_prompt: |
      审查此 Pull Request：
      ${{ github.event.pull_request.body }}
```

### 从文件加载输入提示词

您也可以从文件加载输入提示词：

```yaml
- name: Analyze Code from File
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    model: "gpt-4"
    system_prompt: "你是一位代码分析员"
    input_prompt: "src/main.go" # 从文件加载代码
```

### 从 URL 加载输入提示词

从远程 URL 加载输入内容：

```yaml
- name: Analyze Remote Content
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    system_prompt: "你是一位内容分析员"
    input_prompt: "https://raw.githubusercontent.com/user/repo/main/content.txt"
```

### 在提示词中使用 Go 模板

`system_prompt` 和 `input_prompt` 都支持 Go 模板语法，让您可以动态地将环境变量插入到提示词中。这在 GitHub Actions 工作流程中特别有用，可以包含仓库名称、分支名称或自定义变量等上下文信息。

**主要功能：**

- 使用 `{{.VAR_NAME}}` 访问任何环境变量
- 带有 `INPUT_` 前缀的环境变量可以使用有或没有前缀的形式访问
  - 例如：`INPUT_MODEL` 可以用 `{{.MODEL}}` 或 `{{.INPUT_MODEL}}` 访问
- 所有 GitHub Actions 默认环境变量都可使用（例如 `GITHUB_REPOSITORY`、`GITHUB_REF_NAME`）
- 支持完整的 Go 模板语法，包括条件语句和函数

#### 范例 1：使用 GitHub Actions 变量

```yaml
- name: Analyze Repository with Context
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    model: "gpt-4o"
    system_prompt: |
      你是一位专业的代码分析师。
      专注于 {{.GITHUB_REPOSITORY}} 仓库的分析。
    input_prompt: |
      请分析此仓库：{{.GITHUB_REPOSITORY}}
      当前分支：{{.GITHUB_REF_NAME}}
      使用模型：{{.MODEL}}

      请提供有关代码质量和潜在改进的见解。
```

#### 范例 2：使用自定义环境变量

```yaml
- name: Set Custom Variables
  run: |
    echo "INPUT_PROJECT_TYPE=web-application" >> $GITHUB_ENV
    echo "INPUT_LANGUAGE=Go" >> $GITHUB_ENV

- name: Code Review with Custom Context
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    system_prompt: |
      你正在审查一个使用 {{.LANGUAGE}} 编写的 {{.PROJECT_TYPE}}。
      专注于 {{.LANGUAGE}} 开发的最佳实践。
    input_prompt: |
      审查 {{.GITHUB_REPOSITORY}} 中的代码变更。
      项目类型：{{.PROJECT_TYPE}}
      编程语言：{{.LANGUAGE}}
```

#### 范例 3：模板文件

创建模板文件 `.github/prompts/review-template.txt`：

```text
请审查 {{.GITHUB_REPOSITORY}} 的 Pull Request。

仓库：{{.GITHUB_REPOSITORY}}
分支：{{.GITHUB_REF_NAME}}
执行者：{{.GITHUB_ACTOR}}
模型：{{.MODEL}}

重点关注：
- 代码质量
- 安全性问题
- 性能影响
```

然后在工作流程中使用：

```yaml
- name: Code Review with Template File
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    model: "gpt-4"
    input_prompt: ".github/prompts/review-template.txt"
```

#### 范例 4：条件逻辑

```yaml
- name: Conditional Prompt
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    input_prompt: |
      分析 {{.GITHUB_REPOSITORY}}
      {{if .DEBUG}}
      启用详细输出和详细说明。
      {{else}}
      提供简洁的摘要。
      {{end}}
```

#### 可用的 GitHub Actions 环境变量

可在模板中使用的常见变量：

- `{{.GITHUB_REPOSITORY}}` - 仓库名称（例如 `owner/repo`）
- `{{.GITHUB_REF_NAME}}` - 分支或标签名称
- `{{.GITHUB_ACTOR}}` - 触发工作流程的用户名称
- `{{.GITHUB_SHA}}` - Commit SHA
- `{{.GITHUB_EVENT_NAME}}` - 触发工作流程的事件
- `{{.GITHUB_WORKFLOW}}` - 工作流程名称
- `{{.GITHUB_RUN_ID}}` - 唯一的工作流程执行 ID
- `{{.GITHUB_RUN_NUMBER}}` - 唯一的工作流程执行编号
- 以及工作流程中可用的任何其他环境变量

### 使用 Tool Schema 的结构化输出

使用 `tool_schema` 通过函数调用从 LLM 获取结构化 JSON 输出。当您需要 LLM 以特定格式返回数据，以便在后续工作流程步骤中轻松解析和使用时，这非常有用。

#### 基本结构化输出

```yaml
- name: Extract City Information
  id: extract
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    input_prompt: "法国的首都是什么？"
    tool_schema: |
      {
        "name": "get_city_info",
        "description": "获取城市信息",
        "parameters": {
          "type": "object",
          "properties": {
            "city": {
              "type": "string",
              "description": "城市名称"
            },
            "country": {
              "type": "string",
              "description": "城市所在国家"
            }
          },
          "required": ["city", "country"]
        }
      }

- name: Use Extracted Data
  run: |
    echo "城市：${{ steps.extract.outputs.city }}"
    echo "国家：${{ steps.extract.outputs.country }}"
```

#### 结构化代码审查

```yaml
- name: Structured Code Review
  id: review
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    model: "gpt-4"
    system_prompt: "你是一位专业的代码审查员。"
    input_prompt: |
      审查此代码：
      ```python
      def divide(a, b):
          return a / b
      ```
    tool_schema: |
      {
        "name": "code_review",
        "description": "结构化代码审查结果",
        "parameters": {
          "type": "object",
          "properties": {
            "score": {
              "type": "integer",
              "description": "代码质量评分 1-10"
            },
            "issues": {
              "type": "array",
              "items": { "type": "string" },
              "description": "发现的问题列表"
            },
            "suggestions": {
              "type": "array",
              "items": { "type": "string" },
              "description": "改进建议列表"
            }
          },
          "required": ["score", "issues", "suggestions"]
        }
      }

- name: Process Review Results
  run: |
    echo "评分：${{ steps.review.outputs.score }}"
    echo "问题：${{ steps.review.outputs.issues }}"
    echo "建议：${{ steps.review.outputs.suggestions }}"
```

#### 从文件加载 Tool Schema

将 schema 存放在文件中以便重复使用：

```yaml
- name: Analyze with Schema File
  id: analyze
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    input_prompt: "分析这段文字的情感：我非常喜欢这个产品！"
    tool_schema: ".github/schemas/sentiment-analysis.json"
```

#### Tool Schema 搭配 Go 模板

在 schema 中使用 Go 模板进行动态配置：

```yaml
- name: Dynamic Schema
  uses: appleboy/LLM-action@v1
  env:
    INPUT_FUNCTION_NAME: analyze_text
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    input_prompt: "分析这段文字"
    tool_schema: |
      {
        "name": "{{.FUNCTION_NAME}}",
        "description": "分析文字内容",
        "parameters": {
          "type": "object",
          "properties": {
            "result": { "type": "string" }
          }
        }
      }
```

### 自托管 / 本地 LLM

```yaml
- name: Call Local LLM
  id: local_llm
  uses: appleboy/LLM-action@v1
  with:
    base_url: "http://localhost:8080/v1"
    api_key: "your-local-api-key"
    model: "llama2"
    skip_ssl_verify: "true"
    input_prompt: "用简单的术语解释量子计算"
```

### 使用自定义 CA 证书

对于使用自签名证书的自托管服务，您可以提供自定义 CA 证书。`ca_cert` 输入支持三种格式：

#### 证书内容

```yaml
- name: Call LLM with CA Certificate Content
  uses: appleboy/LLM-action@v1
  with:
    base_url: "https://your-llm-server.local/v1"
    api_key: ${{ secrets.LLM_API_KEY }}
    ca_cert: |
      -----BEGIN CERTIFICATE-----
      MIIDxTCCAq2gAwIBAgIQAqx...
      -----END CERTIFICATE-----
    input_prompt: "Hello, world!"
```

#### 从文件加载证书

```yaml
- name: Call LLM with CA Certificate File
  uses: appleboy/LLM-action@v1
  with:
    base_url: "https://your-llm-server.local/v1"
    api_key: ${{ secrets.LLM_API_KEY }}
    ca_cert: "/path/to/ca-cert.pem"
    input_prompt: "Hello, world!"
```

或使用 `file://` 前缀：

```yaml
- name: Call LLM with CA Certificate File URI
  uses: appleboy/LLM-action@v1
  with:
    base_url: "https://your-llm-server.local/v1"
    api_key: ${{ secrets.LLM_API_KEY }}
    ca_cert: "file:///path/to/ca-cert.pem"
    input_prompt: "Hello, world!"
```

#### 从 URL 加载证书

```yaml
- name: Call LLM with CA Certificate from URL
  uses: appleboy/LLM-action@v1
  with:
    base_url: "https://your-llm-server.local/v1"
    api_key: ${{ secrets.LLM_API_KEY }}
    ca_cert: "https://your-server.com/ca-cert.pem"
    input_prompt: "Hello, world!"
```

### 搭配 Ollama 使用

```yaml
- name: Call Ollama
  id: ollama
  uses: appleboy/LLM-action@v1
  with:
    base_url: "http://localhost:11434/v1"
    api_key: "ollama"
    model: "llama3"
    system_prompt: "你是一个乐于助人的助手"
    input_prompt: "写一首关于编程的俳句"
```

### 链接多个 LLM 调用

```yaml
- name: Generate Story
  id: generate
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    input_prompt: "写一个关于机器人的短篇故事"
    max_tokens: "500"

- name: Translate Story
  id: translate
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    system_prompt: "你是一位翻译员"
    input_prompt: |
      将以下文字翻译成西班牙文：
      ${{ steps.generate.outputs.response }}

- name: Display Results
  run: |
    echo "原始故事："
    echo "${{ steps.generate.outputs.response }}"
    echo ""
    echo "翻译后的故事："
    echo "${{ steps.translate.outputs.response }}"
```

### 调试模式

启用调试模式以排除问题并检查所有参数：

```yaml
- name: Call LLM with Debug
  id: llm_debug
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    model: "gpt-4"
    system_prompt: "你是一个乐于助人的助手"
    input_prompt: "解释 GitHub Actions 如何运作"
    temperature: "0.8"
    max_tokens: "1500"
    debug: true # 启用调试模式
```

**调试输出范例：**

```txt
=== Debug Mode: All Parameters ===
main.Config{
    BaseURL: "https://api.openai.com/v1",
    APIKey: "sk-ab****xyz9",  // 为了安全而屏蔽
    Model: "gpt-4",
    SkipSSLVerify: false,
    SystemPrompt: "你是一个乐于助人的助手",
    InputPrompt: "解释 GitHub Actions 如何运作",
    Temperature: 0.8,
    MaxTokens: 1500,
    Debug: true
}
===================================
=== Debug Mode: Messages ===
[... 消息详情 ...]
============================
```

**安全说明：** 当启用调试模式时，API 密钥会自动屏蔽（仅显示前 4 个和后 4 个字符），以防止在日志中意外泄露。

## 支持的服务

此 Action 适用于任何 OpenAI 兼容的 API，包括：

- **OpenAI** - `https://api.openai.com/v1`
- **Azure OpenAI** - `https://{your-resource}.openai.azure.com/openai/deployments/{deployment-id}`
- **Ollama** - `http://localhost:11434/v1`
- **LocalAI** - `http://localhost:8080/v1`
- **LM Studio** - `http://localhost:1234/v1`
- **Jan** - `http://localhost:1337/v1`
- **vLLM** - 您的 vLLM 服务器端点
- **Text Generation WebUI** - 您的 WebUI 端点
- 任何其他 OpenAI 兼容的服务

## 安全考量

- 请务必使用 GitHub Secrets 存储 API 密钥：`${{ secrets.YOUR_API_KEY }}`
- 仅在信任的本地/内部服务中使用 `skip_ssl_verify: 'true'`
- 请谨慎处理提示词中的敏感数据，因为它们将被发送到 LLM 服务

## 授权

MIT License - 详见 LICENSE 文件

## 贡献

欢迎贡献！请随时提交 Pull Request。
