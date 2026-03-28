## Context

掃描器透過 `filepath.WalkDir()` 逐一處理檔案，目前所有輸出在掃描完成後才產生。使用者在執行期間看不到任何回饋，無法判斷程式是否卡住或正在處理哪個檔案。

## Goals / Non-Goals

**Goals:**
- 掃描期間於 stderr 即時顯示當前處理的檔案路徑
- 顯示已掃描的檔案計數（格式：`[123] /path/to/file`）
- 每次更新覆寫同一行（使用 `\r`），不產生多餘輸出
- 掃描結束後清除進度列，再輸出正式報告
- json output 模式下不顯示進度（避免污染 stdout/stderr pipeline）

**Non-Goals:**
- 百分比進度條（事先無法知道總檔案數）
- ETA 預估
- 修改現有的 text/json 輸出格式
- 在非 TTY 環境（如 CI pipe）強制顯示進度

## Decisions

### 1. Callback 而非 channel

在 `Options` 加入 `OnProgress func(path string, count int64)` callback，由 scanner 在每個檔案處理後呼叫。

**理由**：callback 最簡單，不需要額外 goroutine 或同步機制。scanner 是單執行緒，callback 在同一 goroutine 呼叫，無 race condition。channel 方案需要額外 goroutine 消費，增加複雜度但無顯著好處。

### 2. 輸出至 stderr，格式 `\r[N] <path>`

使用 `\r`（carriage return）覆寫當前行，讓進度顯示停留在同一行。路徑截斷至終端機寬度（預設 80 字元）避免換行。

**理由**：stderr 與 stdout 分離，不影響 json/text 輸出的 pipe 使用。`\r` 覆寫比 ANSI escape codes 更簡單且相容性更好。

### 3. 掃描結束後輸出空行清除進度列

`main.go` 在取得 `ScanReport` 後、呼叫 `reporter.Write()` 前，輸出 `\r\033[K`（清除當前行）。

**理由**：確保正式報告從乾淨的行開始，不殘留進度文字。

### 4. 僅在 text 模式且 stderr 為 TTY 時啟用

`main.go` 使用 `isatty` 或簡單的 `os.Stderr.Stat()` 判斷是否為 TTY，僅在 text 模式下掛接 callback。

**理由**：json 模式通常用於 pipeline，stderr 進度會干擾自動化工具。非 TTY 環境（如 CI）不需要進度顯示。

## Risks / Trade-offs

- `\r` 在 Windows cmd 可能行為不同 → 此為 Linux/macOS CLI 工具，可接受
- 進度顯示頻率與每個檔案一次掛鉤，超大目錄（百萬檔案）可能有輕微 I/O overhead → callback 極輕量（一次 write syscall），可接受
- 若終端機寬度小於截斷長度，路徑可能顯示不全 → 為預期行為

## Migration Plan

無 breaking change。`OnProgress` 預設為 `nil`，scanner 內部以 `if opts.OnProgress != nil` guard，現有呼叫者不受影響。
