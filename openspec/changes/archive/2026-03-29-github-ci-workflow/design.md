## Context

專案使用 Pants 2.30.0 + Go 1.25.0，目前無任何 CI 設定。
GitHub Actions 是最直接的選擇，因為程式碼預計托管於 GitHub，無需額外服務授權。

Pants 在 CI 環境的特殊考量：
- Pants 需要 Python 環境（用於 PEX launcher），GitHub 的 `ubuntu-latest` 已內建
- Pants 第一次執行會下載 PEX bootstrap，可透過快取 `~/.cache/pants/setup` 加速
- `.pants.d/` 是每次建置的工作目錄，應快取以加速增量建置
- Go module cache（`~/go/pkg/mod`）也應快取，避免重複下載 `gopkg.in/yaml.v3` 等相依

## Goals / Non-Goals

**Goals:**
- 每次 push 至 `main` 及所有 PR 時自動執行 `pants test ::` 與 `pants package cmd/avengine:`
- 快取建置產物以縮短 CI 執行時間
- 測試失敗或建置失敗時讓 CI 標記為紅燈（non-zero exit code）

**Non-Goals:**
- 不自動發布二進位（無 release pipeline）
- 不執行 linting（`golint`、`staticcheck`）— 可作為後續 change
- 不設定 branch protection rules（需在 GitHub UI 操作）

## Decisions

**決策 1：直接呼叫系統 `pants` 而非下載腳本**

CI 環境使用 `pip install pantsbuild.pants==2.30.0` 安裝，或直接使用 `curl` 下載官方 setup script，讓 `pants` 指令可用。選擇官方 bootstrap script（`https://static.pantsbuild.org/setup/pants`），與本機開發環境一致，並透過 `pants.toml` 中的 `pants_version = "2.30.0"` 確保版本固定。

替代方案：使用 `pantsbuild/actions/init-pants` 官方 Action → 功能較完整但引入外部相依；選擇 bootstrap script 更輕量且透明。

**決策 2：快取策略**

以下路徑加入 GitHub Actions cache：
- `~/.cache/pants/setup`：Pants PEX bootstrap（版本固定，可長期快取）
- `~/.pants.d`：Pants 工作目錄快取（以 `pants_version` + `go.sum` hash 為 cache key）
- `~/go/pkg/mod`：Go module cache（以 `go.sum` hash 為 cache key）

**決策 3：workflow 觸發條件**

```yaml
on:
  push:
    branches: [main]
  pull_request:
```

只對 `main` push 觸發，PR 對所有分支觸發，符合一般開源專案慣例。

## Risks / Trade-offs

- [風險] Pants 在 CI 第一次執行時需下載 PEX（~40 MB）→ **緩解**：快取 `~/.cache/pants/setup`
- [風險] Go 1.25.0 在 `ubuntu-latest` 的 `actions/setup-go` 可能不是預設版本 → **緩解**：明確指定 `go-version: "1.25.0"`
- [風險] `.pants.d/` 快取 key 過於寬鬆導致快取污染 → **緩解**：以 `pants_version` + `hashFiles('go.sum')` 組合為 key
