## Why

專案目前缺乏持續整合（CI）流程，每次推送程式碼時無法自動驗證建置與測試是否通過。
建立 GitHub Actions workflow 可在每次 push / pull request 時自動執行 `pants test ::` 與 `pants package cmd/avengine:`，確保主線程式碼始終處於可建置、可測試的狀態。

## What Changes

- 新增 `.github/workflows/ci.yml`，定義 CI pipeline，包含：
  - 觸發條件：push 至 `main` 分支、所有 pull request
  - 安裝 Go 1.25.0 環境
  - 安裝系統層級的 Pants（使用 `pants_version = "2.30.0"`）
  - 執行 `pants test ::`，輸出測試結果
  - 執行 `pants package cmd/avengine:`，驗證二進位可成功建置
  - 快取 Pants 建置快取（`.pants.d/`）以加速後續執行

## Capabilities

### New Capabilities

- `github-ci`：定義 GitHub Actions CI workflow 的觸發條件、執行步驟與環境需求

### Modified Capabilities

（無）

## Impact

- **新增檔案**：`.github/workflows/ci.yml`
- **無程式碼變更**：不影響任何 Go 原始碼或建置設定
- **執行環境**：`ubuntu-latest` runner，Go 1.25.0，Pants 2.30.0
