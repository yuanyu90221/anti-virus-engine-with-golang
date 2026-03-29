### Requirement: 掃描器公開進度回呼
`Options` 結構應包含 `OnProgress func(path string, count int64)` 欄位。掃描器應在每個檔案處理完成後呼叫此回呼，傳入檔案路徑與目前已處理的累計檔案數。若 `OnProgress` 為 nil，掃描器的行為應與原有實作完全相同。

#### Scenario: 每個檔案均觸發回呼
- **WHEN** 設定了 `Options.OnProgress` 且掃描執行於 N 個檔案時
- **THEN** 回呼被呼叫恰好 N 次，每個已處理檔案呼叫一次，累計數從 1 開始單調遞增

#### Scenario: nil 回呼不發生 panic
- **WHEN** `Options.OnProgress` 為 nil 時
- **THEN** 掃描正常完成且不產生錯誤，並輸出與原有實作相同的報告

### Requirement: 掃描器提供檔案計數工具函式
`scanner` 套件應公開 `CountFiles(opts Options) (int64, error)` 函式，走訪 `opts.Dir` 並回傳 `Scan` 實際會處理的檔案數量（套用相同的 FollowLinks 與 MaxFileSizeB 過濾條件）。僅在 `opts.Dir` 不存在或無法存取時回傳錯誤。

#### Scenario: 計數與可掃描檔案數吻合
- **WHEN** 對含有 N 個符合條件檔案的目錄呼叫 `CountFiles` 時
- **THEN** 回傳 N（與 `Scan` 報告中的 TotalFiles 一致，不含錯誤與跳過的檔案）

#### Scenario: 空目錄回傳零
- **WHEN** 對空目錄呼叫 `CountFiles` 時
- **THEN** 回傳 0 且不產生錯誤

### Requirement: 文字模式下於 stderr 顯示掃描進度
在文字輸出模式下，CLI 應於掃描期間在 stderr 顯示單行進度指示。指示器應顯示目前已掃描的檔案數目，以及正在處理的檔案路徑（截斷至 80 字元以內）。每次更新應使用歸位符（`\r`）覆蓋前一行。掃描完成後，進度行應在報告輸出至 stdout 前清除。

#### Scenario: 進度更新顯示檔案數目與路徑
- **WHEN** 在文字模式且 stderr 為 TTY 時進行掃描
- **THEN** 每個檔案更新寫入 `\r[<count>] <截斷路徑>` 至 stderr

#### Scenario: 報告輸出前清除進度行
- **WHEN** 掃描在文字模式下完成
- **THEN** 報告寫入 stdout 前，stderr 收到清除行的序列

### Requirement: JSON 模式或非 TTY 時不顯示進度
當輸出格式為 `json` 或 stderr 非 TTY 時，CLI 不應在 stderr 寫入任何進度輸出。

#### Scenario: JSON 模式抑制進度顯示
- **WHEN** 指定 `--output json` 時
- **THEN** stderr 不出現任何進度輸出

#### Scenario: 非 TTY 抑制進度顯示
- **WHEN** stderr 非 TTY（例如導向至檔案）時
- **THEN** stderr 不出現任何進度輸出
