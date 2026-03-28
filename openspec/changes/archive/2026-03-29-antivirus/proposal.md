## Why

目前缺乏一個輕量、可自訂的本機防毒掃描工具，可整合至 CI/CD 流程或本機開發環境中。本計畫以 Go 實作基於 SHA256 hash 比對的防毒引擎，並使用 Pants Build 管理 monorepo，提供可擴充的特徵資料庫架構。

## What Changes

- 新增 CLI 工具 `avengine`，支援子命令 `scan`
- 新增 SHA256 檔案雜湊計算模組（`internal/hasher`）
- 新增特徵資料庫模組（`internal/sigdb`），使用 YAML 格式、可抽換 Loader 介面
- 新增目錄遞迴掃描模組（`internal/scanner`），支援大小過濾、符號連結控制
- 新增結果輸出模組（`internal/reporter`），支援文字表格與 JSON 格式
- 新增範例特徵資料庫（`signatures/ransomware.yaml`、`signatures/trojans.yaml`）
- 使用 Pants Build 管理整個 monorepo 建構與測試

## Capabilities

### New Capabilities

- `file-hashing`：對單一檔案進行串流 SHA256 計算，回傳十六進位字串
- `signature-db`：從 YAML 目錄載入特徵，建立記憶體索引，提供 hash 查詢（Loader 介面可抽換）
- `directory-scan`：遞迴走訪目錄，對每個檔案計算 hash 並查詢特徵資料庫，產生掃描報告
- `scan-report`：將掃描報告以文字表格或 JSON 格式輸出，並定義結束碼語意
- `avengine-cli`：整合所有模組的 CLI 進入點，支援 `scan` 子命令與各項參數

### Modified Capabilities

（無現有規格需修改）

## Impact

- **新增外部相依**：`gopkg.in/yaml.v3`（唯一外部套件）
- **建構工具**：需安裝 Pants 2.23.0、Go 1.22+
- **新增檔案**：`pants.toml`、`go.mod`、`go.sum`、所有 `internal/` 套件、`cmd/avengine/main.go`、`signatures/*.yaml`、`testdata/`
- **無破壞性變更**（全新專案）
