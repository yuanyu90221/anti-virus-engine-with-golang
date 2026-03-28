## 1. 啟動 Pants

- [x] 1.1 下載 Pants 啟動腳本：`curl -L -o ./pants https://static.pantsbuild.org/setup/pants && chmod +x ./pants`
- [x] 1.2 在儲存庫根目錄建立 `pants.toml`，設定 `pants_version = "2.23.0"`、`backend_packages = ["pants.backend.experimental.go"]`、`[source] root_patterns = ["/"]`，以及 `[golang] go_search_paths = ["<GOROOT>/bin", "<PATH>"]`
- [x] 1.3 驗證 `./pants --version` 輸出為 `2.23.0`

## 2. 初始化 Go 模組

- [x] 2.1 執行 `go mod init github.com/yuanyu/avengine` 建立 `go.mod`
- [x] 2.2 驗證 `go mod verify` 以代碼 0 結束

## 3. 建立根目錄 BUILD 檔案

- [x] 3.1 在儲存庫根目錄建立空白 `BUILD` 檔案（僅含註解即可；`pants tailor` 將自動填充內容）

## 4. 驗證

- [x] 4.1 執行 `./pants help goals`，確認輸出中包含 `test`、`package` 與 `tailor` 目標
- [x] 4.2 建立最小化佔位檔案 `cmd/avengine/main.go`，內容為 `package main` 與 `func main() {}`
- [x] 4.3 執行 `./pants tailor ::`，確認以代碼 0 結束並產生 BUILD 條目
- [x] 4.4 視情況保留或刪除佔位 `main.go`（若後續立即開始實作，則保留）
