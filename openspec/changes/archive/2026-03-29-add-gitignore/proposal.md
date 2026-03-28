## Why

專案目前的 `.gitignore` 僅涵蓋標準 Go 二進位檔與測試產物，缺少 Pants 建置快取（`.pants.d/`）、建置輸出（`dist/`）、以及執行時產生的掃描報告（`scan-report.json`）等項目，導致這些檔案可能被意外提交至版本控制。

## What Changes

- 在現有 `.gitignore` 中新增以下類別的忽略規則：
  - **Pants 建置快取與工作目錄**：`.pants.d/`、`.pids`
  - **建置輸出目錄**：`dist/`
  - **執行時產生的掃描報告**：`scan-report.json`、`*.scan.json`
  - **IDE 設定目錄**（補齊原本被注解的規則）：`.idea/`、`.vscode/`
  - **作業系統暫存檔**：`.DS_Store`、`Thumbs.db`

## Capabilities

### New Capabilities

- `gitignore-config`：定義專案 `.gitignore` 應涵蓋的忽略規則類別與具體項目

### Modified Capabilities

（無）— 不涉及任何對外 API 或業務邏輯的規格變更

## Impact

- **`.gitignore`**：新增規則，不影響現有已追蹤的檔案
- **版本控制**：防止建置產物、快取、IDE 設定被意外提交
- **開發者體驗**：`git status` 輸出更乾淨，`git add .` 更安全
