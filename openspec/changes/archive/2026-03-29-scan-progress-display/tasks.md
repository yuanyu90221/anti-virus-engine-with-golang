## 1. Scanner Progress Callback

- [x] 1.1 在 `internal/scanner/scanner.go` 的 `Options` struct 新增 `OnProgress func(path string, count int64)` 欄位
- [x] 1.2 在 `Scan()` 的 `WalkDir` loop 中，每處理完一個檔案後呼叫 `opts.OnProgress`（guard `!= nil`）

## 2. CLI Progress Display

- [x] 2.1 在 `cmd/avengine/main.go` 新增 `isTerminal(f *os.File) bool` 輔助函式（使用 `os.File.Stat()` 判斷 TTY）
- [x] 2.2 在 text 模式且 stderr 為 TTY 時，建立 progress callback：寫入 `\r[N] <path>` 至 stderr（路徑截斷至 80 字元）
- [x] 2.3 掃描結束後、呼叫 `reporter.Write()` 前，清除 stderr 進度列（寫入 `\r\033[K`）
- [x] 2.4 json 模式或非 TTY 時不設定 `OnProgress`（保持 nil）

## 3. 驗證

- [x] 3.1 本地執行 `pants test ::` 確認現有測試全數通過
- [x] 3.2 手動執行 `go run ./cmd/avengine scan --dir ./testdata` 確認進度列正確顯示並在報告前清除
- [x] 3.3 執行 `go run ./cmd/avengine scan --dir ./testdata --output json` 確認 json 模式無進度輸出
