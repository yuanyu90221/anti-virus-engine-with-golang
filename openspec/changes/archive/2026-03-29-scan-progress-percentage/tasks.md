## 1. Scanner CountFiles 函式

- [x] 1.1 在 `internal/scanner/scanner.go` 新增 `CountFiles(opts Options) (int64, error)` 函式，走訪 `opts.Dir` 並套用相同的 FollowLinks / MaxFileSizeB 過濾邏輯，回傳符合條件的檔案總數

## 2. CLI 百分比進度顯示

- [x] 2.1 在 `cmd/avengine/main.go` 的進度啟用區段中，於建立 callback 前呼叫 `scanner.CountFiles(scanOpts)` 取得 `total`
- [x] 2.2 當 `total > 0` 時，進度 closure 格式改為 `\r[N/Total] (XX%) path`（路徑截斷至 70 字元）
- [x] 2.3 當 `total == 0` 或 CountFiles 回傳 error 時，退回原本 `\r[N] path` 格式

## 3. 驗證

- [x] 3.1 執行 `pants test ::` 確認所有測試通過
- [x] 3.2 執行 `go run ./cmd/avengine scan --dir ./testdata --sigs ./signatures` 確認進度列顯示 `[N/Total] (XX%)` 格式
- [x] 3.3 執行 `go run ./cmd/avengine scan --dir ./testdata --sigs ./signatures --output json` 確認 json 模式無進度輸出
