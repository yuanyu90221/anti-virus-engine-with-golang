## Context

專案根目錄已存在 `.gitignore`，內容為標準 Go 模板（二進位、測試產物、coverage），但缺少以下在本專案實際產生的檔案：

- **`.pants.d/`**：Pants 建置快取目錄，含 `pids/`、`workdir/`、`bin/` 等子目錄，每次建置都會更新
- **`dist/`**：`pants package` 輸出的二進位檔（如 `dist/cmd.avengine/bin`）
- **`scan-report.json`**：README 的 CI/CD 範例腳本執行後會產生此檔
- **`.idea/`、`.vscode/`**：IDE 設定（原本被注解，應啟用）

## Goals / Non-Goals

**Goals:**
- 確保所有建置產物、快取、IDE 設定、執行時輸出不被意外追蹤
- 直接修改現有 `.gitignore`，不新增檔案

**Non-Goals:**
- 不使用 `.gitignore` 全域樣板（global gitignore）
- 不處理 `openspec/` 目錄下的任何內容（這些應被追蹤）
- 不移除現有規則

## Decisions

**決策：直接附加至現有 `.gitignore`，依類別分組**

以注解區塊清楚標示新增的規則來源，讓未來維護者能理解每條規則的用途。不使用單一大型規則替換整個檔案，保留原有 Go 標準模板的完整性。

## Risks / Trade-offs

- [風險] `dist/` 中若有需要追蹤的資產（如預建二進位發布版本）→ **緩解**：本專案以原始碼建置為主，`dist/` 應視為建置輸出而非版本控制資產
- [風險] 已被追蹤的 `scan-report.json` 若存在於 git index 中，`.gitignore` 不會自動取消追蹤 → **緩解**：此為新增規則，本專案目前無 git history，無此問題
