## Why

目前進度列只顯示已掃描的檔案數（`[N] path`），使用者無法判斷還剩多少工作量。加入百分比顯示（`[N/Total] (XX%) path`）可讓使用者清楚掌握整體進度。

## What Changes

- 在 `scanner` 套件新增 `CountFiles(opts Options) (int64, error)` 函式，對目標目錄做一次快速前置走訪，計算符合條件（非略過）的檔案總數
- 在 `cmd/avengine/main.go` 的進度顯示邏輯中，啟用進度時先呼叫 `CountFiles` 取得總數，再將百分比加入進度列格式：`[N/Total] (XX%) path`
- `ProgressFunc` 簽章不變，`total` 透過 closure 傳遞，無 breaking change

## Capabilities

### New Capabilities

（無）

### Modified Capabilities

- `scan-progress`：進度列格式從 `[N] path` 改為 `[N/Total] (XX%) path`，需加入前置計數步驟

## Impact

- `internal/scanner/scanner.go`：新增 `CountFiles` 函式
- `cmd/avengine/main.go`：在 progress 模式下先呼叫 `CountFiles`，更新 closure 格式
- 不影響現有測試（`CountFiles` 為新增函式）
- 不影響 json 模式或非 TTY 環境
