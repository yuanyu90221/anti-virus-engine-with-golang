## Context

專案使用 Pants 2.30.0 + Go 1.25.0，CI workflow 已建立。
Release workflow 需要在 tag 觸發時建置二進位並發布至 GitHub Releases。

GitHub 官方提供 `softprops/action-gh-release`（社群廣泛採用）作為建立 Release 的 Action，
可在單一步驟內同時建立 Release 並上傳 assets，比舊版的 `actions/create-release` + `actions/upload-release-asset` 組合更簡潔。

## Goals / Non-Goals

**Goals:**
- 推送 `v*` tag 後自動建立 GitHub Release 並附上 `avengine` 二進位
- Release 名稱與 tag 名稱一致（如 tag `v1.0.0` → Release `v1.0.0`）
- 二進位重命名為 `avengine`（而非 Pants 預設的 `bin`）

**Non-Goals:**
- 不支援多平台交叉編譯（Go 1.25 + Pants 目前僅建置當前 runner 平台，linux/amd64）
- 不自動產生 changelog 或 release notes
- 不建置 `.tar.gz` 壓縮包（僅裸二進位）

## Decisions

**決策 1：使用 `softprops/action-gh-release@v2`**

單一步驟完成建立 Release + 上傳 assets，配置簡單。
需要 `permissions: contents: write` 讓 `GITHUB_TOKEN` 有寫入 Releases 的權限。

替代方案：`gh release create`（GitHub CLI）→ 需要額外安裝步驟，較繁瑣。

**決策 2：二進位重命名**

Pants 輸出路徑為 `dist/cmd.avengine/bin`，上傳前以 `mv` 重命名為 `avengine`，
讓下載者得到直觀的執行檔名稱，不需再次重命名。

**決策 3：快取策略與 CI 一致**

複用相同的 Go module cache 與 Pants cache 設定，確保建置環境一致。

## Risks / Trade-offs

- [風險] `pants package` 在 CI 環境建置的二進位為 linux/amd64，不支援 macOS/Windows → **緩解**：文件說明平台限制，未來可用 matrix 擴充
- [風險] `GITHUB_TOKEN` 預設可能缺少 `contents: write` 權限（取決於 repo 設定） → **緩解**：在 workflow 中明確宣告 `permissions`
