## Why

目前掃描執行時沒有任何視覺回饋，面對大型目錄時使用者無法判斷程式是否仍在運行。加入即時進度顯示可讓使用者清楚知道掃描進度與當前處理的檔案。

## What Changes

- 在 `Scanner.Scan()` 加入 progress callback 機制，每處理一個檔案即回呼通知
- 在 `cmd/avengine/main.go` 連接 progress callback，於 stderr 輸出當前掃描檔案名稱與進度條
- 進度條顯示：已掃描檔案數 / 格式化當前檔案路徑（截斷過長路徑）
- 掃描結束後清除進度列，再輸出正式報告至 stdout（不干擾 text/json 輸出）

## Capabilities

### New Capabilities
- `scan-progress`: 掃描時在 stderr 即時顯示當前處理檔案與進度計數

### Modified Capabilities
- （無 spec-level 行為變更）

## Impact

- `internal/scanner/scanner.go`：新增 `ProgressFunc` 型別與 `Options.OnProgress` 欄位
- `cmd/avengine/main.go`：在 text output 模式下連接進度顯示邏輯（json 模式不顯示）
- 不影響現有輸出格式（progress 輸出至 stderr，報告仍輸出至 stdout）
- 不影響現有測試（callback 為 optional，預設為 nil）
