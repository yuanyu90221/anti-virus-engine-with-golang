## 為何需要此變更

本專案需要 Pants 建置環境，以啟用 Go 模組管理、增量測試與二進位封裝功能。若沒有 `pants.toml` 與有效的 `go.mod`，任何建置或測試指令皆無法執行。

## 變更內容

- 新增 `pants.toml`，設定 Pants 2.23.0 並啟用 Go 後端
- 新增根目錄 `BUILD` 檔案（最小化內容；`pants tailor ::` 將自動填充）
- 初始化模組 `github.com/yuanyu/avengine` 的 `go.mod`
- 將 `pants` 啟動腳本加入專案根目錄

## 功能模組

### 新增功能模組

- `pants-go-setup`：設定 Pants 2.23.0 並啟用 `pants.backend.go`、初始化 `go.mod`，以及建立最小化根目錄 BUILD 檔案，使 `./pants --version`、`./pants tailor ::` 與 `./pants test ::` 能正常運作。

### 修改現有功能模組

## 影響範圍

- 新增 `pants.toml`、根目錄 `BUILD`、`go.mod` 與 `pants` 啟動腳本作為受版控的新檔案
- 不影響任何現有 Go 原始碼（目前尚無原始碼）
- 相依條件：主機上必須已安裝 Go 1.22 以上版本
