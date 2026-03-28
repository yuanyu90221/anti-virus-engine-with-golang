## ADDED Requirements

### Requirement: Release workflow 在推送 v* tag 時自動觸發
`.github/workflows/release.yml` 中 SHALL 設定在推送符合 `v*` 格式的 tag 時觸發 release pipeline。

#### Scenario: 推送語意化版本 tag 時觸發
- **WHEN** 開發者執行 `git push origin v1.0.0`
- **THEN** GitHub Actions 自動啟動 release workflow

#### Scenario: 一般 branch push 不觸發 release workflow
- **WHEN** 開發者推送 commit 至 `main` 分支（非 tag）
- **THEN** release workflow 不執行（僅 CI workflow 執行）

### Requirement: Release workflow 建置 avengine 二進位
Release workflow SHALL 執行 `pants package cmd/avengine:` 並取得建置產物。

#### Scenario: 二進位建置成功
- **WHEN** release workflow 在 `ubuntu-latest` 上執行 `pants package cmd/avengine:`
- **THEN** `dist/cmd.avengine/bin` 存在且為可執行的 ELF 二進位

### Requirement: Release workflow 建立 GitHub Release 並以 tag 名稱為版本號
Release workflow SHALL 自動建立 GitHub Release，Release 名稱與 tag 名稱一致。

#### Scenario: Release 名稱與 tag 一致
- **WHEN** 推送 tag `v1.2.3` 後 release workflow 完成
- **THEN** GitHub Releases 頁面出現名稱為 `v1.2.3` 的 Release

### Requirement: Release artifact 重命名為 avengine 並可供下載
Release workflow SHALL 將建置產物重命名為 `avengine` 後上傳至 GitHub Release assets。

#### Scenario: 下載連結指向 avengine 檔案
- **WHEN** 使用者在 GitHub Releases 頁面查看 Release assets
- **THEN** 可下載名稱為 `avengine` 的二進位檔案

### Requirement: Release workflow 具備寫入 Releases 的權限
`.github/workflows/release.yml` 中 SHALL 明確宣告 `permissions: contents: write`，
確保 `GITHUB_TOKEN` 可建立 Release 與上傳 assets。

#### Scenario: GITHUB_TOKEN 可建立 Release
- **WHEN** release workflow 呼叫 GitHub Releases API
- **THEN** 不因權限不足而失敗（HTTP 403）
