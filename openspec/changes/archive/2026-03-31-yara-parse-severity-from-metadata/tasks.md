## 1. 更新 CLI 呼叫與輸出解析

- [x] 1.1 在 `internal/yara/yara.go` 的 CLI 呼叫中加入 `--print-meta` 旗標
- [x] 1.2 新增 `parseYaraLine(line string) (name, severity string)` 輔助函數，支援含 `[...]` 與不含 `[...]` 兩種輸出格式
- [x] 1.3 新增 `parseSeverityFromMeta(meta string) string` 輔助函數，從 `key="val",...` 格式中擷取 `severity` 值
- [x] 1.4 更新 `parseOutput()` 改為呼叫 `parseYaraLine()`，移除原本的 `strings.SplitN` 邏輯

## 2. 更新測試

- [x] 2.1 更新 `TestInspect_Match` 的 fake binary 輸出為新格式（含 metadata block），並將 severity 斷言從 `"unknown"` 改為實際值
- [x] 2.2 更新 `TestInspect_MultipleMatches` 的 fake binary 輸出，加入各自的 severity metadata
- [x] 2.3 新增測試案例 `TestInspect_NoSeverityMeta`：規則無 `severity` metadata 時，`Severity` 應 fallback 為 `"unknown"`
- [x] 2.4 新增單元測試直接測試 `parseYaraLine()` 與 `parseSeverityFromMeta()`（涵蓋有值、無值、多欄位等情境）

## 3. 驗證

- [x] 3.1 執行 `go test ./internal/yara/...` 確認全數通過
- [x] 3.2 執行 `go test ./...` 確認無回歸
