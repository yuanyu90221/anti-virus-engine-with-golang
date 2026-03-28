## Requirements

### Requirement: 文字表格格式輸出掃描摘要與偵測結果
系統 SHALL 在 `format = "text"` 時，將 `ScanReport` 以人類可讀的文字表格格式寫入 `io.Writer`，包含：掃描目錄、檔案總數、威脅數、錯誤數、略過數、耗時，以及每筆偵測結果（路徑、SHA256 前 16 字元、威脅名稱、嚴重程度、分類）。

#### Scenario: 無威脅時的文字輸出
- **WHEN** `ScanReport.Detections` 為空且 format 為 "text"
- **THEN** 輸出包含「未發現威脅」或類似提示，且不包含偵測結果表格列

#### Scenario: 有威脅時的文字輸出
- **WHEN** `ScanReport.Detections` 不為空且 format 為 "text"
- **THEN** 輸出包含每筆偵測的路徑、SHA256 前 16 字元、威脅名稱、嚴重程度與分類

### Requirement: JSON 格式輸出完整掃描報告
系統 SHALL 在 `format = "json"` 時，將完整 `ScanReport` 序列化為 JSON 並寫入 `io.Writer`，欄位名稱使用 camelCase。

#### Scenario: JSON 輸出可被解析
- **WHEN** format 為 "json" 且掃描完成
- **THEN** 輸出為合法 JSON，可被 `json.Unmarshal` 解析，且包含 `detections`、`totalFiles`、`errorFiles`、`skippedFiles` 等欄位

### Requirement: 定義三種結束碼語意
系統 SHALL 定義以下常數：`ExitClean = 0`（無威脅）、`ExitDetected = 1`（發現威脅）、`ExitError = 2`（致命錯誤），並由 CLI 層據此呼叫 `os.Exit`。

#### Scenario: 掃描完成且無威脅
- **WHEN** 掃描完成且 `Detections` 為空
- **THEN** CLI 以結束碼 0 退出

#### Scenario: 掃描完成且發現威脅
- **WHEN** 掃描完成且 `Detections` 不為空
- **THEN** CLI 以結束碼 1 退出

#### Scenario: 致命錯誤（例如特徵目錄不存在）
- **WHEN** 無法載入特徵資料庫或其他致命錯誤發生
- **THEN** CLI 以結束碼 2 退出並輸出錯誤訊息至 stderr
