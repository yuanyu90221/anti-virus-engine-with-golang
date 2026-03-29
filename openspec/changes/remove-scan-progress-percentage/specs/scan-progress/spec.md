## MODIFIED Requirements

### Requirement: 文字模式下於 stderr 顯示掃描進度
在文字輸出模式下，CLI 應於掃描期間在 stderr 顯示單行進度指示。指示器應顯示目前已掃描的檔案數目，以及正在處理的檔案路徑（截斷至 80 字元以內）。每次更新應使用歸位符（`\r`）覆蓋前一行。掃描完成後，進度行應在報告輸出至 stdout 前清除。

#### Scenario: 進度更新顯示檔案數目與路徑
- **WHEN** 在文字模式且 stderr 為 TTY 時進行掃描
- **THEN** 每個檔案更新寫入 `\r[<count>] <截斷路徑>` 至 stderr

#### Scenario: 報告輸出前清除進度行
- **WHEN** 掃描在文字模式下完成
- **THEN** 報告寫入 stdout 前，stderr 收到清除行的序列

## REMOVED Requirements

### Requirement: CLI 掃描時使用 CountFiles 取得檔案總數
**Reason**: CLI 不再於掃描前呼叫 `CountFiles`。移除預掃描步驟可消除啟動延遲，並移除百分比顯示。`CountFiles` 函式仍保留於 `scanner` 套件，但不再由 CLI 進度路徑呼叫。
**Migration**: 無需任何操作；`CountFiles` 仍可供程式呼叫使用。
