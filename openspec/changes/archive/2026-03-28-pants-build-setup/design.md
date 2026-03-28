## 背景

全新的 Go 專案，目前尚無任何建置設定。Pants 2.23.0 已安裝於 `/home/yuanyu/.local/bin/pants`。本專案將使用單一 Go 模組（`github.com/yuanyu/avengine`）搭配 `pants.backend.experimental.go` 後端，以實現增量建置與測試。

## 目標 / 非目標

**目標：**
- 建立 `pants.toml`，啟用 `pants.backend.experimental.go` 並設定版本為 2.23.0
- 為 `github.com/yuanyu/avengine` 初始化 `go.mod`
- 提供最小化根目錄 `BUILD` 檔案；`pants tailor ::` 將自動產生其餘內容
- 驗證 `./pants --version` 與 `./pants tailor ::` 可無誤執行

**非目標：**
- 撰寫任何 Go 原始碼
- 新增外部 Go 相依套件（`go get ...`）——此部分由防毒引擎實作變更負責
- 建立 CI/CD 流水線
- 設定 `pants.ci.toml` 或 lint/format 相關組態

## 決策

| 決策項目 | 選擇 | 理由 |
|----------|------|------|
| Pants 版本 | 2.23.0 | 與現有防毒引擎設計規格一致；為支援 Go 後端的穩定發行版 |
| Go 搜尋路徑 | `["<GOROOT>/bin", "<PATH>"]` | 讓 Pants 自動探索 PATH 上的 Go 安裝位置，避免硬編碼路徑 |
| 原始碼根目錄 | `["/"]` | 單一儲存庫結構——所有套件皆位於儲存庫根目錄之下 |
| 模組名稱 | `github.com/yuanyu/avengine` | 與現有 design.md 的決策保持一致 |

`pants.toml` 內容：
```toml
[GLOBAL]
pants_version = "2.23.0"
backend_packages = ["pants.backend.experimental.go"]

[source]
root_patterns = ["/"]

[golang]
go_search_paths = ["<GOROOT>/bin", "<PATH>"]
```

## 風險 / 取捨

- **Go 版本不符** → Pants 2.23.0 需要 Go 1.21 以上；設定前請先執行 `go version` 確認。
- **需要網路連線** → 首次執行 `./pants` 時會下載 Pants PEX 二進位檔；需有網路存取。
- **`pants tailor` 會產生 BUILD 檔案** → 請勿手動編輯自動產生的 BUILD 檔案；新增套件後重新執行 `tailor`。

## 部署計畫

1. 將 `pants.toml` 與根目錄 `BUILD` 加入儲存庫根目錄
2. 執行 `go mod init github.com/yuanyu/avengine`
3. 執行 `./pants --version` 驗證啟動是否正常
4. 待 Go 原始碼存在後，執行 `./pants tailor ::`

回滾方式：刪除 `pants.toml`、`BUILD`、`go.mod` 與 `.pants.d/` 快取目錄。
