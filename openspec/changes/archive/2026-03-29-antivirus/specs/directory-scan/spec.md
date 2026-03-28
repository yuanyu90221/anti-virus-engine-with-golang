## ADDED Requirements

### Requirement: 遞迴掃描目錄中的每個檔案
系統 SHALL 使用 `fs.WalkDir` 遞迴走訪 `Options.Dir`，對每個一般檔案計算 SHA256 並查詢 `sigdb.DB`，最終產生 `ScanReport`。

#### Scenario: 乾淨目錄掃描
- **WHEN** 目錄中所有檔案的 hash 均不在特徵資料庫中
- **THEN** `ScanReport.Detections` 為空 slice，`TotalFiles` 等於掃描的檔案數，`ErrorFiles` 與 `SkippedFiles` 為 0

#### Scenario: 包含惡意檔案的目錄掃描
- **WHEN** 目錄中至少一個檔案的 hash 命中特徵資料庫
- **THEN** 該檔案出現在 `ScanReport.Detections` 中，含路徑、hash、MatchResult 資訊

### Requirement: 預設不追蹤符號連結
系統 SHALL 在 `Options.FollowLinks = false`（預設）時略過符號連結，將其計入 `SkippedFiles`。僅當 `Options.FollowLinks = true` 時才追蹤。

#### Scenario: 目錄含符號連結且未啟用 follow-links
- **WHEN** `FollowLinks` 為 false 且目錄中存在符號連結
- **THEN** 符號連結不被掃描，`SkippedFiles` 計數加一

### Requirement: 超過大小限制的檔案應被略過
系統 SHALL 在 `Options.MaxFileSizeB > 0` 時，略過大小超過該值（bytes）的檔案，並將其計入 `SkippedFiles`。`MaxFileSizeB = 0` 表示不限制。

#### Scenario: 檔案超過大小限制
- **WHEN** `MaxFileSizeB` 為 1048576（1 MB）且目錄中存在 2 MB 的檔案
- **THEN** 該檔案不被雜湊計算，`SkippedFiles` 加一

### Requirement: 無法讀取的檔案應被記錄為錯誤
系統 SHALL 在無法讀取某檔案時，將其計入 `ErrorFiles`，並繼續掃描其餘檔案（不中止整個掃描）。

#### Scenario: 無讀取權限的檔案
- **WHEN** 目錄中存在一個無讀取權限的檔案
- **THEN** `ErrorFiles` 加一，掃描繼續，其他檔案正常處理

### Requirement: ScanReport 記錄掃描時間
系統 SHALL 在掃描開始時記錄 `StartedAt`，掃描結束時記錄 `FinishedAt`，兩者均為 `time.Time`。

#### Scenario: 掃描完成後時間戳記有效
- **WHEN** `Scan()` 正常完成
- **THEN** `FinishedAt` 晚於或等於 `StartedAt`，且兩者均不為零值
