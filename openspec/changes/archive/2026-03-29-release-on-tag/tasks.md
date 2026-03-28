## 1. 建立 release workflow 檔案

- [x] 1.1 建立 `.github/workflows/release.yml` 檔案

## 2. 設定觸發條件與權限

- [x] 2.1 設定 `on.push.tags: ['v*']`，僅在推送 v* tag 時觸發
- [x] 2.2 設定 `permissions: contents: write`，允許 GITHUB_TOKEN 建立 Release

## 3. 設定執行環境

- [x] 3.1 設定 job 使用 `ubuntu-latest` runner
- [x] 3.2 加入 `actions/checkout@v4` 步驟（含 `fetch-depth: 0` 以取得完整 git history）
- [x] 3.3 加入 `actions/setup-go@v5` 步驟，指定 `go-version: "1.25"`

## 4. 設定快取與 Pants 安裝

- [x] 4.1 加入 Go module cache（`~/go/pkg/mod`）快取步驟
- [x] 4.2 加入 Pants cache（`~/.cache/pants`）快取步驟
- [x] 4.3 加入下載 Pants bootstrap script 的步驟

## 5. 設定建置與發布步驟

- [x] 5.1 加入執行 `pants package cmd/avengine:` 的步驟
- [x] 5.2 加入將 `dist/cmd.avengine/bin` 複製並重命名為 `avengine` 的步驟
- [x] 5.3 加入 `softprops/action-gh-release@v2` 步驟，上傳 `avengine` 並建立 Release

## 6. 驗證

- [x] 6.1 確認 `.github/workflows/release.yml` YAML 語法正確
- [x] 6.2 確認 tag 觸發條件（`v*`）與權限宣告（`contents: write`）存在
- [x] 6.3 確認 artifact 路徑與重命名邏輯正確（`dist/cmd.avengine/bin` → `avengine`）
