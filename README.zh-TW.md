# LLM Action

[English](README.md) | [繁體中文](README.zh-TW.md) | [簡體中文](README.zh-CN.md)

[![Lint and Testing](https://github.com/appleboy/LLM-action/actions/workflows/testing.yml/badge.svg)](https://github.com/appleboy/LLM-action/actions/workflows/testing.yml)
[![Trivy Security Scan](https://github.com/appleboy/LLM-action/actions/workflows/trivy.yml/badge.svg)](https://github.com/appleboy/LLM-action/actions/workflows/trivy.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/LLM-action)](https://goreportcard.com/report/github.com/appleboy/LLM-action)

一個用於與 OpenAI 相容 LLM 服務互動的 GitHub Action。此 Action 允許您連接到任何 OpenAI 相容的 API 端點（包括本地或自架服務），並獲取可用於工作流程的回應。

## 功能特色

- 🔌 連接任何 OpenAI 相容的 API 端點
- 🔐 支援自訂 API 金鑰
- 🔧 可配置的基礎 URL，適用於自架服務
- 🚫 選擇性跳過 SSL 憑證驗證
- 🎯 支援系統提示詞以設定情境
- 📝 輸出回應可用於後續 Actions
- 🎛️ 可配置的溫度和最大權杖數
- 🐛 偵錯模式，並安全地遮罩 API 金鑰
- 🎨 支援 Go 模板語法，可動態插入環境變數

## 輸入參數

| 輸入              | 說明                                                                              | 必填 | 預設值                      |
| ----------------- | --------------------------------------------------------------------------------- | ---- | --------------------------- |
| `base_url`        | OpenAI 相容 API 端點的基礎 URL                                                    | 否   | `https://api.openai.com/v1` |
| `api_key`         | 用於驗證的 API 金鑰                                                               | 是   | -                           |
| `model`           | 要使用的模型名稱                                                                  | 否   | `gpt-4o`                    |
| `skip_ssl_verify` | 跳過 SSL 憑證驗證                                                                 | 否   | `false`                     |
| `system_prompt`   | 設定情境的系統提示詞。支援純文字、檔案路徑或 URL。支援 Go 模板語法與環境變數      | 否   | `''`                        |
| `input_prompt`    | 使用者輸入給 LLM 的提示詞。支援純文字、檔案路徑或 URL。支援 Go 模板語法與環境變數 | 是   | -                           |
| `temperature`     | 回應隨機性的溫度值（0.0-2.0）                                                     | 否   | `0.7`                       |
| `max_tokens`      | 回應中的最大權杖數                                                                | 否   | `1000`                      |
| `debug`           | 啟用偵錯模式以顯示所有參數（API 金鑰將被遮罩）                                    | 否   | `false`                     |

## 輸出參數

| 輸出       | 說明            |
| ---------- | --------------- |
| `response` | 來自 LLM 的回應 |

## 使用範例

### 基本範例

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

### 使用系統提示詞

````yaml
- name: Code Review with LLM
  id: review
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    model: "gpt-4"
    system_prompt: "你是一位程式碼審查員。請提供有關程式碼品質、最佳實務和潛在問題的建設性意見。"
    input_prompt: |
      請審查此程式碼：
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

### 使用多行系統提示詞

```yaml
- name: Advanced Code Review
  id: review
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    model: "gpt-4"
    system_prompt: |
      你是一位擁有深厚軟體工程最佳實務知識的專業程式碼審查員。

      你的職責：
      - 識別潛在的錯誤和安全漏洞
      - 建議改善程式碼品質和可維護性的方法
      - 檢查是否遵守程式碼標準
      - 評估效能影響

      請以專業的語氣提供建設性、可行的意見。
    input_prompt: |
      審查以下 Pull Request 變更：
      ${{ github.event.pull_request.body }}
    temperature: "0.3"
    max_tokens: "2000"
```

### 從檔案載入系統提示詞

不需要在 YAML 中嵌入冗長的提示詞，可以從檔案載入：

````yaml
- name: Code Review with Prompt File
  id: review
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    model: "gpt-4"
    system_prompt: ".github/prompts/code-review.txt"
    input_prompt: |
      審查此程式碼：
      ```python
      def calculate(x, y):
          return x / y
      ```
````

或使用 `file://` 前綴：

```yaml
- name: Code Review with File URI
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    system_prompt: "file://.github/prompts/code-review.txt"
    input_prompt: "審查 main.go 檔案"
```

### 從 URL 載入系統提示詞

從遠端 URL 載入提示詞：

```yaml
- name: Code Review with Remote Prompt
  id: review
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    model: "gpt-4"
    system_prompt: "https://raw.githubusercontent.com/your-org/prompts/main/code-review.txt"
    input_prompt: |
      審查此 Pull Request：
      ${{ github.event.pull_request.body }}
```

### 從檔案載入輸入提示詞

您也可以從檔案載入輸入提示詞：

```yaml
- name: Analyze Code from File
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    model: "gpt-4"
    system_prompt: "你是一位程式碼分析員"
    input_prompt: "src/main.go" # 從檔案載入程式碼
```

### 從 URL 載入輸入提示詞

從遠端 URL 載入輸入內容：

```yaml
- name: Analyze Remote Content
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    system_prompt: "你是一位內容分析員"
    input_prompt: "https://raw.githubusercontent.com/user/repo/main/content.txt"
```

### 在提示詞中使用 Go 模板

`system_prompt` 和 `input_prompt` 都支援 Go 模板語法，讓您可以動態地將環境變數插入到提示詞中。這在 GitHub Actions 工作流程中特別有用，可以包含儲存庫名稱、分支名稱或自訂變數等上下文資訊。

**主要功能：**

- 使用 `{{.VAR_NAME}}` 存取任何環境變數
- 帶有 `INPUT_` 前綴的環境變數可以使用有或沒有前綴的形式存取
  - 例如：`INPUT_MODEL` 可以用 `{{.MODEL}}` 或 `{{.INPUT_MODEL}}` 存取
- 所有 GitHub Actions 預設環境變數都可使用（例如 `GITHUB_REPOSITORY`、`GITHUB_REF_NAME`）
- 支援完整的 Go 模板語法，包括條件式和函數

#### 範例 1：使用 GitHub Actions 變數

```yaml
- name: Analyze Repository with Context
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    model: "gpt-4o"
    system_prompt: |
      你是一位專業的程式碼分析師。
      專注於 {{.GITHUB_REPOSITORY}} 儲存庫的分析。
    input_prompt: |
      請分析此儲存庫：{{.GITHUB_REPOSITORY}}
      目前分支：{{.GITHUB_REF_NAME}}
      使用模型：{{.MODEL}}

      請提供有關程式碼品質和潛在改進的見解。
```

#### 範例 2：使用自訂環境變數

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
      你正在審查一個使用 {{.LANGUAGE}} 撰寫的 {{.PROJECT_TYPE}}。
      專注於 {{.LANGUAGE}} 開發的最佳實務。
    input_prompt: |
      審查 {{.GITHUB_REPOSITORY}} 中的程式碼變更。
      專案類型：{{.PROJECT_TYPE}}
      程式語言：{{.LANGUAGE}}
```

#### 範例 3：模板檔案

建立模板檔案 `.github/prompts/review-template.txt`：

```text
請審查 {{.GITHUB_REPOSITORY}} 的 Pull Request。

儲存庫：{{.GITHUB_REPOSITORY}}
分支：{{.GITHUB_REF_NAME}}
執行者：{{.GITHUB_ACTOR}}
模型：{{.MODEL}}

重點關注：
- 程式碼品質
- 安全性問題
- 效能影響
```

然後在工作流程中使用：

```yaml
- name: Code Review with Template File
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    model: "gpt-4"
    input_prompt: ".github/prompts/review-template.txt"
```

#### 範例 4：條件邏輯

```yaml
- name: Conditional Prompt
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    input_prompt: |
      分析 {{.GITHUB_REPOSITORY}}
      {{if .DEBUG}}
      啟用詳細輸出和詳細說明。
      {{else}}
      提供簡潔的摘要。
      {{end}}
```

#### 可用的 GitHub Actions 環境變數

可在模板中使用的常見變數：

- `{{.GITHUB_REPOSITORY}}` - 儲存庫名稱（例如 `owner/repo`）
- `{{.GITHUB_REF_NAME}}` - 分支或標籤名稱
- `{{.GITHUB_ACTOR}}` - 觸發工作流程的使用者名稱
- `{{.GITHUB_SHA}}` - Commit SHA
- `{{.GITHUB_EVENT_NAME}}` - 觸發工作流程的事件
- `{{.GITHUB_WORKFLOW}}` - 工作流程名稱
- `{{.GITHUB_RUN_ID}}` - 唯一的工作流程執行 ID
- `{{.GITHUB_RUN_NUMBER}}` - 唯一的工作流程執行編號
- 以及工作流程中可用的任何其他環境變數

### 自架 / 本地 LLM

```yaml
- name: Call Local LLM
  id: local_llm
  uses: appleboy/LLM-action@v1
  with:
    base_url: "http://localhost:8080/v1"
    api_key: "your-local-api-key"
    model: "llama2"
    skip_ssl_verify: "true"
    input_prompt: "用簡單的術語解釋量子計算"
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
    system_prompt: "你是一個樂於助人的助手"
    input_prompt: "寫一首關於程式設計的俳句"
```

### 鏈結多個 LLM 呼叫

```yaml
- name: Generate Story
  id: generate
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    input_prompt: "寫一個關於機器人的短篇故事"
    max_tokens: "500"

- name: Translate Story
  id: translate
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    system_prompt: "你是一位翻譯員"
    input_prompt: |
      將以下文字翻譯成西班牙文：
      ${{ steps.generate.outputs.response }}

- name: Display Results
  run: |
    echo "原始故事："
    echo "${{ steps.generate.outputs.response }}"
    echo ""
    echo "翻譯後的故事："
    echo "${{ steps.translate.outputs.response }}"
```

### 偵錯模式

啟用偵錯模式以排除問題並檢查所有參數：

```yaml
- name: Call LLM with Debug
  id: llm_debug
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    model: "gpt-4"
    system_prompt: "你是一個樂於助人的助手"
    input_prompt: "解釋 GitHub Actions 如何運作"
    temperature: "0.8"
    max_tokens: "1500"
    debug: true # 啟用偵錯模式
```

**偵錯輸出範例：**

```txt
=== Debug Mode: All Parameters ===
main.Config{
    BaseURL: "https://api.openai.com/v1",
    APIKey: "sk-ab****xyz9",  // 為了安全而遮罩
    Model: "gpt-4",
    SkipSSLVerify: false,
    SystemPrompt: "你是一個樂於助人的助手",
    InputPrompt: "解釋 GitHub Actions 如何運作",
    Temperature: 0.8,
    MaxTokens: 1500,
    Debug: true
}
===================================
=== Debug Mode: Messages ===
[... 訊息詳情 ...]
============================
```

**安全說明：** 當啟用偵錯模式時，API 金鑰會自動遮罩（僅顯示前 4 個和後 4 個字元），以防止在日誌中意外洩露。

## 支援的服務

此 Action 適用於任何 OpenAI 相容的 API，包括：

- **OpenAI** - `https://api.openai.com/v1`
- **Azure OpenAI** - `https://{your-resource}.openai.azure.com/openai/deployments/{deployment-id}`
- **Ollama** - `http://localhost:11434/v1`
- **LocalAI** - `http://localhost:8080/v1`
- **LM Studio** - `http://localhost:1234/v1`
- **Jan** - `http://localhost:1337/v1`
- **vLLM** - 您的 vLLM 伺服器端點
- **Text Generation WebUI** - 您的 WebUI 端點
- 任何其他 OpenAI 相容的服務

## 安全考量

- 請務必使用 GitHub Secrets 儲存 API 金鑰：`${{ secrets.YOUR_API_KEY }}`
- 僅在信任的本地/內部服務中使用 `skip_ssl_verify: 'true'`
- 請謹慎處理提示詞中的敏感資料，因為它們將被發送到 LLM 服務

## 授權

MIT License - 詳見 LICENSE 文件

## 貢獻

歡迎貢獻！請隨時提交 Pull Request。
