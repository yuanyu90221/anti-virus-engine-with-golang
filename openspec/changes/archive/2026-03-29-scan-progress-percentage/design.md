## Context

現有進度列（`scan-progress-display` change）已實作 `[N] path` 格式，透過 `OnProgress` callback 在每個檔案處理後回呼。百分比需要事先知道總檔案數，但 `filepath.WalkDir` 是單次走訪，無法在走訪過程中得知總數。

## Goals / Non-Goals

**Goals:**
- 掃描前快速預計數（pre-count），取得符合條件的檔案總數
- 進度列格式升級為 `[N/Total] (XX%) path`
- 不改變 `ProgressFunc` 型別簽章

**Non-Goals:**
- 修改 json 模式或非 TTY 行為
- 對預計數結果做快取或持久化

## Decisions

### 1. 新增 `CountFiles(opts Options) (int64, error)` 函式

在 `scanner` 套件新增獨立的前置計數函式，走訪目錄並套用相同的過濾條件（FollowLinks、MaxFileSizeB），回傳符合條件的檔案總數。

**理由**：將計數邏輯封裝在 scanner 套件，與 `Scan` 共用相同的過濾邏輯，避免在 main.go 重複實作。獨立函式比修改 `Scan` 介面更簡單，且可被測試。

**替代方案**：在 `Options` 加入 `TotalFiles int64` 欄位由外部傳入 → 需呼叫者自行計數，邏輯外漏，不採用。

### 2. `total` 透過 closure 傳遞給 `OnProgress`

`ProgressFunc` 簽章 `func(path string, count int64)` 保持不變。`main.go` 在建立 callback 前先呼叫 `CountFiles`，再在 closure 中捕捉 `total` 變數：

```go
total, _ := scanner.CountFiles(opts)
onProgress = func(path string, count int64) {
    pct := count * 100 / total
    fmt.Fprintf(os.Stderr, "\r[%d/%d] (%d%%) %s", count, total, pct, display)
}
```

**理由**：零 breaking change，不影響任何現有呼叫者。

### 3. CountFiles 遇到 total=0 時的安全處理

若 `CountFiles` 回傳 0（空目錄）或錯誤，`main.go` 退回原本的 `[N] path` 格式，避免除以零。

## Risks / Trade-offs

- 預計數增加一次額外走訪，IO overhead 與目錄規模成正比 → 對掃描大目錄（百萬檔案）略有影響，但計數走訪不做 hash，速度遠快於完整掃描，可接受
- 預計數與實際掃描之間若有檔案異動，百分比可能超過 100% 或提早達到 100% → 為預期行為，不影響正確性
