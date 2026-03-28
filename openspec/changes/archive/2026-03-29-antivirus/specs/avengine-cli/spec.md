## ADDED Requirements

### Requirement: scan 子命令接受必要與可選參數
系統 SHALL 提供 `avengine scan` 子命令，支援以下旗標：
- `--dir <path>`（必填）：遞迴掃描的目標目錄
- `--sigs <path>`（選填，預設 `./signatures`）：特徵 YAML 目錄
- `--output text|json`（選填，預設 `text`）：輸出格式
- `--follow-links`（選填，預設 false）：追蹤符號連結
- `--max-size <MB>`（選填，預設 0 = 不限制）：略過超過 N MB 的檔案
- `--verbose`（選填，預設 false）：顯示所有檔案，不只偵測到的項目

#### Scenario: 僅提供必填參數 --dir
- **WHEN** 執行 `avengine scan --dir ./testdata`
- **THEN** 工具使用預設值（`./signatures`、`text` 格式、不追蹤連結、不限大小）完成掃描並輸出結果

#### Scenario: 未提供 --dir 參數
- **WHEN** 執行 `avengine scan`（不帶 `--dir`）
- **THEN** 工具輸出使用說明至 stderr 並以結束碼 2 退出

#### Scenario: --output 傳入無效值
- **WHEN** 執行 `avengine scan --dir . --output xml`
- **THEN** 工具輸出錯誤訊息至 stderr 並以結束碼 2 退出

### Requirement: CLI 整合所有內部模組並以正確結束碼退出
系統 SHALL 依序執行：載入特徵（`sigdb.NewDB`）→ 掃描目錄（`scanner.Scan`）→ 輸出報告（`reporter.Write`）→ 依報告結果呼叫 `os.Exit(0|1|2)`。

#### Scenario: 掃描乾淨目錄的完整流程
- **WHEN** 指定目錄中無惡意檔案且特徵資料庫載入成功
- **THEN** 工具輸出掃描摘要並以結束碼 0 退出

#### Scenario: 掃描含惡意檔案的目錄完整流程
- **WHEN** 指定目錄中存在 hash 命中特徵資料庫的檔案
- **THEN** 工具輸出含偵測結果的報告並以結束碼 1 退出

#### Scenario: 特徵目錄不存在
- **WHEN** `--sigs` 指定的目錄不存在
- **THEN** 工具輸出錯誤訊息至 stderr 並以結束碼 2 退出
