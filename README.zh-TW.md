# LLM Action

[English](README.md) | [繁體中文](README.zh-TW.md) | [簡體中文](README.zh-CN.md)

![LLM Action](images/llm-action_800x600.png)

[![Lint and Testing](https://github.com/appleboy/LLM-action/actions/workflows/testing.yml/badge.svg)](https://github.com/appleboy/LLM-action/actions/workflows/testing.yml)
[![Trivy Security Scan](https://github.com/appleboy/LLM-action/actions/workflows/trivy.yml/badge.svg)](https://github.com/appleboy/LLM-action/actions/workflows/trivy.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/LLM-action)](https://goreportcard.com/report/github.com/appleboy/LLM-action)

一個用於與 OpenAI 相容 LLM 服務互動的 GitHub Action，支援自訂端點、自架模型（Ollama、LocalAI、vLLM）、SSL/CA 憑證、Go template 動態提示詞，以及透過 function calling 實現結構化輸出。

## 簡報

了解如何使用此 Action 打造 AI 驅動的 GitHub 自動化工作流程：

- [打造AI驅動的GitHub自動化工作流程](https://speakerdeck.com/appleboy/da-zao-a-i-qu-dong-de-g-i-t-h-u-b-dong-hua-zuo-liu-cheng) - 涵蓋 Tool Schema 結構化輸出、LLM 服務無縫切換，以及實際應用場景如程式碼審查、PR 摘要和 Issue 分類。

## 為什麼選擇 LLM Action？

隨著 AI 輔助開發成為主流，將 LLM 整合到 CI/CD 流程中對現代軟體團隊來說至關重要。然而，現有的解決方案通常面臨以下挑戰：

- **廠商綁定**：大多數 GitHub Actions 只支援單一 LLM 服務商，難以切換服務或比較不同模型
- **彈性不足**：寫死的提示詞和僵化的設定無法適應多樣化的專案需求
- **安全疑慮**：自架或隔離網路環境需要自訂端點和憑證管理，許多 Actions 不支援
- **非結構化輸出**：原始 LLM 回應難以解析，不利於自動化工作流程使用

**LLM Action** 的誕生就是為了解決這些問題：

1. **通用相容性** - 一個 Action 支援所有 OpenAI 相容服務（OpenAI、Azure、Ollama、LocalAI、vLLM 等）
2. **動態提示詞** - Go 模板讓你注入環境變數，打造情境感知的提示詞
3. **結構化輸出** - Tool schema（function calling）確保一致、可解析的 JSON 回應，便於自動化處理
4. **企業級支援** - 完整支援自訂 CA 憑證、SSL 設定和 HTTP headers，適用於安全部署環境

無論你要打造自動化程式碼審查、PR 摘要、Issue 分類，或任何 AI 驅動的工作流程，LLM Action 都能讓你靈活使用任何 LLM 服務，同時保持工作流程的可移植性和可維護性。

## 目錄

- [LLM Action](#llm-action)
  - [簡報](#簡報)
  - [為什麼選擇 LLM Action？](#為什麼選擇-llm-action)
  - [目錄](#目錄)
  - [功能特色](#功能特色)
  - [輸入參數](#輸入參數)
  - [輸出參數](#輸出參數)
  - [使用範例](#使用範例)
    - [基本範例](#基本範例)
    - [版本固定](#版本固定)
    - [使用系統提示詞](#使用系統提示詞)
    - [使用多行系統提示詞](#使用多行系統提示詞)
    - [從檔案載入系統提示詞](#從檔案載入系統提示詞)
    - [從 URL 載入系統提示詞](#從-url-載入系統提示詞)
    - [從檔案載入輸入提示詞](#從檔案載入輸入提示詞)
    - [從 URL 載入輸入提示詞](#從-url-載入輸入提示詞)
    - [在提示詞中使用 Go 模板](#在提示詞中使用-go-模板)
      - [範例 1：使用 GitHub Actions 變數](#範例-1使用-github-actions-變數)
      - [範例 2：使用自訂環境變數](#範例-2使用自訂環境變數)
      - [範例 3：模板檔案](#範例-3模板檔案)
      - [範例 4：條件邏輯](#範例-4條件邏輯)
      - [可用的 GitHub Actions 環境變數](#可用的-github-actions-環境變數)
    - [使用 Tool Schema 的結構化輸出](#使用-tool-schema-的結構化輸出)
      - [基本結構化輸出](#基本結構化輸出)
      - [結構化程式碼審查](#結構化程式碼審查)
      - [從檔案載入 Tool Schema](#從檔案載入-tool-schema)
      - [Tool Schema 搭配 Go 模板](#tool-schema-搭配-go-模板)
      - [處理陣列與巢狀物件](#處理陣列與巢狀物件)
    - [自架 / 本地 LLM](#自架--本地-llm)
    - [搭配 Azure OpenAI 使用](#搭配-azure-openai-使用)
    - [使用自訂 CA 憑證](#使用自訂-ca-憑證)
      - [憑證內容](#憑證內容)
      - [從檔案載入憑證](#從檔案載入憑證)
      - [從 URL 載入憑證](#從-url-載入憑證)
    - [搭配 Ollama 使用](#搭配-ollama-使用)
    - [鏈結多個 LLM 呼叫](#鏈結多個-llm-呼叫)
    - [偵錯模式](#偵錯模式)
    - [自訂 HTTP Headers](#自訂-http-headers)
      - [預設 Headers](#預設-headers)
      - [自訂 Headers](#自訂-headers)
      - [單行格式](#單行格式)
      - [多行格式](#多行格式)
      - [搭配自訂認證使用](#搭配自訂認證使用)
  - [支援的服務](#支援的服務)
  - [安全考量](#安全考量)
  - [授權](#授權)
  - [貢獻](#貢獻)

## 功能特色

- 🔌 連接任何 OpenAI 相容的 API 端點
- 🔐 支援自訂 API 金鑰
- 🔧 可配置的基礎 URL，適用於自架服務
- 🚫 選擇性跳過 SSL 憑證驗證
- 🔒 支援自訂 CA 憑證，適用於自簽憑證
- 🎯 支援系統提示詞以設定情境
- 📝 輸出回應可用於後續 Actions
- 🎛️ 可配置的溫度和最大權杖數
- 🐛 偵錯模式，並安全地遮罩 API 金鑰
- 🎨 支援 Go 模板語法，可動態插入環境變數
- 🛠️ 透過函數呼叫支援結構化輸出（tool schema 支援）
- 📋 支援自訂 HTTP headers，適用於日誌分析和自訂認證

## 輸入參數

| 輸入                      | 說明                                                                                   | 必填 | 預設值                      |
| ----------------------- | -------------------------------------------------------------------------------------- | ---- | --------------------------- |
| `base_url`              | OpenAI 相容 API 端點的基礎 URL                                                         | 否   | `https://api.openai.com/v1` |
| `api_key`               | 用於驗證的 API 金鑰                                                                    | 是   | -                           |
| `model`                 | 要使用的模型名稱                                                                       | 否   | `gpt-4o`                    |
| `skip_ssl_verify`       | 跳過 SSL 憑證驗證                                                                      | 否   | `false`                     |
| `ca_cert`               | 自訂 CA 憑證。支援憑證內容、檔案路徑或 URL                                             | 否   | `''`                        |
| `system_prompt`         | 設定情境的系統提示詞。支援純文字、檔案路徑或 URL。支援 Go 模板語法與環境變數           | 否   | `''`                        |
| `input_prompt`          | 使用者輸入給 LLM 的提示詞。支援純文字、檔案路徑或 URL。支援 Go 模板語法與環境變數      | 是   | -                           |
| `tool_schema`           | 用於結構化輸出的 JSON schema（函數呼叫）。支援純文字、檔案路徑或 URL。支援 Go 模板語法 | 否   | `''`                        |
| `temperature`           | 回應隨機性的溫度值（0.0-2.0）                                                          | 否   | `0.7`                       |
| `max_tokens`            | 回應中的最大權杖數                                                                     | 否   | `1000`                      |
| `max_completion_tokens` | 推理模型（o1/o3/o4/gpt-5 系列）的最大完成權杖數。優先於 `max_tokens` | 否 | `''` |
| `debug`                 | 啟用偵錯模式以顯示所有參數（API 金鑰將被遮罩）                                         | 否   | `false`                     |
| `headers`               | 自訂 HTTP headers。格式：`Header1:Value1,Header2:Value2` 或多行格式                    | 否   | `''`                        |

新設定建議使用 `max_completion_tokens`：

```yaml
with:
  max_completion_tokens: "2000"
```

它優先於 `max_tokens`。舊服務若只支援 `max_tokens`，仍可保留；未設定 `max_completion_tokens` 時，Action 會自動為可辨識的推理模型（o1/o3/o4/gpt-5 系列）將 `max_tokens` 轉為 `max_completion_tokens` 傳送。

## 輸出參數

| 輸出                                    | 說明                                                              |
| --------------------------------------- | ----------------------------------------------------------------- |
| `response`                              | 來自 LLM 的原始回應（始終可用）                                   |
| `prompt_tokens`                         | 提示詞的 token 數量                                               |
| `completion_tokens`                     | 回覆的 token 數量                                                 |
| `total_tokens`                          | 總 token 使用量                                                   |
| `prompt_cached_tokens`                  | 提示詞中的快取 token 數量（節省成本，如可用）                     |
| `completion_reasoning_tokens`           | 推理 token 數量，用於 o1/o3 模型（如可用）                        |
| `completion_accepted_prediction_tokens` | 已接受的預測 token 數量（如可用）                                 |
| `completion_rejected_prediction_tokens` | 已拒絕的預測 token 數量（如可用）                                 |
| `<field>`                               | 使用 tool_schema 時，函數參數 JSON 中的每個欄位都會成為獨立的輸出 |

**輸出行為：**

- `response` 輸出**始終可用**，包含 LLM 的原始回應
- 當使用 `tool_schema` 時，函數參數會被解析，每個欄位都會作為獨立輸出加入，同時保留 `response`
- **保留欄位：** 如果您的 tool schema 定義了名為 `response` 的欄位，該欄位將被**跳過**並顯示警告訊息。這是因為 `response` 是保留給 LLM 原始輸出使用的

**範例：** 如果您的 schema 定義了 `city` 和 `country` 欄位，輸出將會是：

- `steps.<id>.outputs.response` - 原始 JSON 回應
- `steps.<id>.outputs.city` - city 欄位的值
- `steps.<id>.outputs.country` - country 欄位的值

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

### 版本固定

您可以固定此 Action 的特定版本：

```yaml
# 使用主版本（推薦 - 自動獲取相容的更新）
uses: appleboy/LLM-action@v1

# 使用特定版本（最大穩定性）
uses: appleboy/LLM-action@v1.0.0

# 使用最新開發版本（不建議用於生產環境）
uses: appleboy/LLM-action@main
```

**建議：** 使用主版本標籤（例如 `@v1`）以自動獲取向後相容的更新和錯誤修復。

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

### 使用 Tool Schema 的結構化輸出

使用 `tool_schema` 透過函數呼叫從 LLM 獲取結構化 JSON 輸出。當您需要 LLM 以特定格式回傳資料，以便在後續工作流程步驟中輕鬆解析和使用時，這非常有用。

#### 基本結構化輸出

```yaml
- name: Extract City Information
  id: extract
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    input_prompt: "法國的首都是什麼？"
    tool_schema: |
      {
        "name": "get_city_info",
        "description": "取得城市資訊",
        "parameters": {
          "type": "object",
          "properties": {
            "city": {
              "type": "string",
              "description": "城市名稱"
            },
            "country": {
              "type": "string",
              "description": "城市所在國家"
            }
          },
          "required": ["city", "country"]
        }
      }

- name: Use Extracted Data
  run: |
    echo "城市：${{ steps.extract.outputs.city }}"
    echo "國家：${{ steps.extract.outputs.country }}"
```

#### 結構化程式碼審查

````yaml
- name: Structured Code Review
  id: review
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    model: "gpt-4"
    system_prompt: "你是一位專業的程式碼審查員。"
    input_prompt: |
      審查此程式碼：
      ```python
      def divide(a, b):
          return a / b
      ```
    tool_schema: |
      {
        "name": "code_review",
        "description": "結構化程式碼審查結果",
        "parameters": {
          "type": "object",
          "properties": {
            "score": {
              "type": "integer",
              "description": "程式碼品質評分 1-10"
            },
            "issues": {
              "type": "array",
              "items": { "type": "string" },
              "description": "發現的問題列表"
            },
            "suggestions": {
              "type": "array",
              "items": { "type": "string" },
              "description": "改進建議列表"
            }
          },
          "required": ["score", "issues", "suggestions"]
        }
      }

- name: Process Review Results
  env:
    SCORE: ${{ steps.review.outputs.score }}
    ISSUES: ${{ steps.review.outputs.issues }}
    SUGGESTIONS: ${{ steps.review.outputs.suggestions }}
  run: |
    echo "評分：$SCORE"
    echo "問題：$ISSUES"
    echo "建議：$SUGGESTIONS"
````

**為什麼使用環境變數而非直接插值？**

- **自動轉義特殊字符**：GitHub Actions 會自動處理環境變數中的特殊符號，避免 shell 解析錯誤
- **更安全**：防止注入攻擊和意外的命令執行，特別是處理 LLM 輸出時
- **更清晰**：程式碼更易讀且易於維護

#### 從檔案載入 Tool Schema

將 schema 存放在檔案中以便重複使用：

```yaml
- name: Analyze with Schema File
  id: analyze
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    input_prompt: "分析這段文字的情感：我非常喜歡這個產品！"
    tool_schema: ".github/schemas/sentiment-analysis.json"
```

#### Tool Schema 搭配 Go 模板

在 schema 中使用 Go 模板進行動態配置：

```yaml
- name: Dynamic Schema
  uses: appleboy/LLM-action@v1
  env:
    INPUT_FUNCTION_NAME: analyze_text
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    input_prompt: "分析這段文字"
    tool_schema: |
      {
        "name": "{{.FUNCTION_NAME}}",
        "description": "分析文字內容",
        "parameters": {
          "type": "object",
          "properties": {
            "result": { "type": "string" }
          }
        }
      }
```

#### 處理陣列與巢狀物件

GitHub Actions 的輸出永遠是字串。當你的 tool schema 回傳陣列或巢狀物件時，它們會被序列化為 JSON 字串。在後續步驟中使用 GitHub 的 `fromJSON()` 函數來解析它們。

**陣列輸出範例：**

```yaml
- name: 擷取關鍵字
  id: keywords
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    input_prompt: "從以下文字擷取關鍵字：GitHub Actions 自動化 CI/CD 工作流程"
    tool_schema: |
      {
        "name": "extract_keywords",
        "description": "從文字中擷取關鍵字",
        "parameters": {
          "type": "object",
          "properties": {
            "keywords": {
              "type": "array",
              "items": { "type": "string" },
              "description": "擷取的關鍵字列表"
            }
          },
          "required": ["keywords"]
        }
      }

- name: 使用關鍵字陣列
  run: |
    # keywords 輸出是 JSON 字串：["GitHub","Actions","CI/CD","工作流程"]
    # 使用 fromJSON() 來解析
    echo "第一個關鍵字: ${{ fromJSON(steps.keywords.outputs.keywords)[0] }}"
    echo "所有關鍵字: ${{ join(fromJSON(steps.keywords.outputs.keywords), ', ') }}"
```

**巢狀物件範例：**

```yaml
- name: 分析程式碼結構
  id: analysis
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    input_prompt: "分析一個 React 元件的結構"
    tool_schema: |
      {
        "name": "analyze_code",
        "description": "分析程式碼結構",
        "parameters": {
          "type": "object",
          "properties": {
            "component": {
              "type": "object",
              "properties": {
                "name": { "type": "string" },
                "props": {
                  "type": "array",
                  "items": { "type": "string" }
                },
                "hooks": {
                  "type": "array",
                  "items": { "type": "string" }
                }
              }
            }
          },
          "required": ["component"]
        }
      }

- name: 使用巢狀資料
  run: |
    # 使用 fromJSON() 存取巢狀屬性
    echo "元件: ${{ fromJSON(steps.analysis.outputs.component).name }}"
    echo "第一個 prop: ${{ fromJSON(steps.analysis.outputs.component).props[0] }}"
    echo "使用的 hooks: ${{ join(fromJSON(steps.analysis.outputs.component).hooks, ', ') }}"
```

**動態 Matrix 範例：**

使用陣列輸出來建立動態的 job matrix：

```yaml
jobs:
  analyze:
    runs-on: ubuntu-latest
    outputs:
      targets: ${{ steps.llm.outputs.targets }}
    steps:
      - name: 取得建置目標
        id: llm
        uses: appleboy/LLM-action@v1
        with:
          api_key: ${{ secrets.OPENAI_API_KEY }}
          input_prompt: "列出要建置的平台：linux、macos、windows"
          tool_schema: |
            {
              "name": "get_targets",
              "description": "取得建置目標",
              "parameters": {
                "type": "object",
                "properties": {
                  "targets": {
                    "type": "array",
                    "items": { "type": "string" }
                  }
                },
                "required": ["targets"]
              }
            }

  build:
    needs: analyze
    strategy:
      matrix:
        target: ${{ fromJSON(needs.analyze.outputs.targets) }}
    runs-on: ubuntu-latest
    steps:
      - run: echo "正在為 ${{ matrix.target }} 建置"
```

**重要注意事項：**

- 所有非字串值（陣列、物件、數字、布林值）都會被 JSON 序列化為字串
- 使用 `fromJSON()` 將字串轉換回原始類型
- 對於大整數，請注意 JSON 解析中可能的浮點數精度問題
- 深層巢狀結構可能需要多次呼叫 `fromJSON()`

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

### 搭配 Azure OpenAI 使用

Azure OpenAI 服務需要不同的 URL 格式。您需要在基礎 URL 中指定資源名稱和部署 ID：

```yaml
- name: Call Azure OpenAI
  id: azure_llm
  uses: appleboy/LLM-action@v1
  with:
    base_url: "https://{your-resource-name}.openai.azure.com/openai/deployments/{deployment-id}"
    api_key: ${{ secrets.AZURE_OPENAI_API_KEY }}
    model: "gpt-4" # 應與您部署的模型相符
    system_prompt: "你是一個樂於助人的助手"
    input_prompt: "說明雲端運算的優點"
```

**設定說明：**

- 將 `{your-resource-name}` 替換為您的 Azure OpenAI 資源名稱
- 將 `{deployment-id}` 替換為您的模型部署名稱
- `model` 參數應與您部署的模型相符
- API 金鑰可在 Azure Portal 中您的 OpenAI 資源的「金鑰和端點」下找到

**完整參數範例：**

```yaml
- name: Azure OpenAI Code Review
  id: azure_review
  uses: appleboy/LLM-action@v1
  with:
    base_url: "https://my-openai-resource.openai.azure.com/openai/deployments/gpt-4-deployment"
    api_key: ${{ secrets.AZURE_OPENAI_API_KEY }}
    model: "gpt-4"
    system_prompt: "你是一位專業的程式碼審查員"
    input_prompt: |
      審查此程式碼的最佳實務：
      ${{ github.event.pull_request.body }}
    temperature: "0.3"
    max_tokens: "2000"
```

### 使用自訂 CA 憑證

對於使用自簽憑證的自架服務，您可以提供自訂 CA 憑證。`ca_cert` 輸入支援三種格式：

#### 憑證內容

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

#### 從檔案載入憑證

```yaml
- name: Call LLM with CA Certificate File
  uses: appleboy/LLM-action@v1
  with:
    base_url: "https://your-llm-server.local/v1"
    api_key: ${{ secrets.LLM_API_KEY }}
    ca_cert: "/path/to/ca-cert.pem"
    input_prompt: "Hello, world!"
```

或使用 `file://` 前綴：

```yaml
- name: Call LLM with CA Certificate File URI
  uses: appleboy/LLM-action@v1
  with:
    base_url: "https://your-llm-server.local/v1"
    api_key: ${{ secrets.LLM_API_KEY }}
    ca_cert: "file:///path/to/ca-cert.pem"
    input_prompt: "Hello, world!"
```

#### 從 URL 載入憑證

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

### 自訂 HTTP Headers

#### 預設 Headers

每個 API 請求都會自動包含以下 headers，用於識別和日誌分析：

| Header             | 值                     | 說明                              |
| ------------------ | ---------------------- | --------------------------------- |
| `User-Agent`       | `LLM-action/{version}` | 標準 User-Agent，包含 Action 版本 |
| `X-Action-Name`    | `appleboy/LLM-action`  | GitHub Action 的完整名稱          |
| `X-Action-Version` | `{version}`            | Action 的語意化版本號             |

這些 headers 可協助您在 LLM 服務日誌中識別來自此 Action 的請求。

#### 自訂 Headers

使用 `headers` 輸入參數為 API 請求添加自訂 HTTP headers。適用於：

- 添加請求追蹤 ID 以進行日誌分析
- 自訂認證標頭
- 傳遞元資料給您的 LLM 服務

#### 單行格式

```yaml
- name: Call LLM with Custom Headers
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    input_prompt: "Hello, world!"
    headers: "X-Request-ID:${{ github.run_id }},X-Trace-ID:${{ github.sha }}"
```

#### 多行格式

```yaml
- name: Call LLM with Multiple Headers
  uses: appleboy/LLM-action@v1
  with:
    api_key: ${{ secrets.OPENAI_API_KEY }}
    input_prompt: "分析此程式碼"
    headers: |
      X-Request-ID:${{ github.run_id }}
      X-Trace-ID:${{ github.sha }}
      X-Environment:production
      X-Repository:${{ github.repository }}
```

#### 搭配自訂認證使用

```yaml
- name: Call Custom LLM Service
  uses: appleboy/LLM-action@v1
  with:
    base_url: "https://your-llm-service.com/v1"
    api_key: ${{ secrets.LLM_API_KEY }}
    input_prompt: "產生摘要"
    headers: |
      X-Custom-Auth:${{ secrets.CUSTOM_AUTH_TOKEN }}
      X-Tenant-ID:my-tenant
```

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
