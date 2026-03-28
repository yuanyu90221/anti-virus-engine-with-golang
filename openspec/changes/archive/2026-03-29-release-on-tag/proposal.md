## Why

目前 CI workflow 只在 push/PR 時驗證建置與測試，缺少版本發布流程。
當開發者推送 git tag（如 `v1.0.0`）時，需要能自動以該 tag 作為版本號建置二進位檔，
並上傳至 GitHub Releases 供使用者直接下載，取代手動建置與上傳的繁瑣步驟。

## What Changes

- 新增 `.github/workflows/release.yml`，定義 release pipeline：
  - 觸發條件：push 符合 `v*` 格式的 tag（如 `v1.0.0`、`v0.2.1`）
  - 安裝與 CI 相同的 Go 1.25 + Pants 2.30.0 環境
  - 執行 `pants package cmd/avengine:` 建置二進位
  - 自動建立 GitHub Release，以 tag 名稱為版本號
  - 將建置產物（`dist/cmd.avengine/bin`）重新命名為 `avengine` 並上傳至 Release

## Capabilities

### New Capabilities

- `release-on-tag`：定義 GitHub Actions release workflow 的觸發條件、建置流程與 artifact 上傳規格

### Modified Capabilities

（無）

## Impact

- **新增檔案**：`.github/workflows/release.yml`
- **無程式碼變更**：不影響任何 Go 原始碼或建置設定
- **GitHub Releases**：每次推送 `v*` tag 後自動產生 Release 頁面與下載連結
