## 1. 建立 workflow 目錄與檔案結構

- [x] 1.1 建立 `.github/workflows/` 目錄
- [x] 1.2 建立 `.github/workflows/ci.yml` 檔案

## 2. 設定 workflow 觸發條件

- [x] 2.1 設定 `on.push.branches: [main]`，push 至 main 時觸發
- [x] 2.2 設定 `on.pull_request`，所有 PR 時觸發

## 3. 設定執行環境

- [x] 3.1 設定 job 使用 `ubuntu-latest` runner
- [x] 3.2 加入 `actions/checkout@v4` 步驟以取得程式碼
- [x] 3.3 加入 `actions/setup-go@v5` 步驟，指定 `go-version: "1.25"`

## 4. 設定 Pants 安裝步驟

- [x] 4.1 加入下載 Pants bootstrap script 的步驟（`curl -fsSL https://static.pantsbuild.org/setup/pants -o ./pants && chmod +x ./pants`）

## 5. 設定快取

- [x] 5.1 加入 `actions/cache@v4` 快取 Go module cache（`~/go/pkg/mod`），以 `go.sum` hash 為 key
- [x] 5.2 加入 `actions/cache@v4` 快取 Pants 建置快取（`~/.cache/pants`），以 `pants_version` + `go.sum` hash 為 key

## 6. 設定建置與測試步驟

- [x] 6.1 加入執行 `pants test ::` 的步驟，命名為「Run tests」
- [x] 6.2 加入執行 `pants package cmd/avengine:` 的步驟，命名為「Build binary」

## 7. 驗證

- [x] 7.1 確認 `.github/workflows/ci.yml` 語法正確（可用 `actionlint` 或直接推送至 GitHub 驗證）
- [x] 7.2 確認 workflow 中的 Go 版本（`1.25`）與 `go.mod` 的 `go 1.25.0` 一致
