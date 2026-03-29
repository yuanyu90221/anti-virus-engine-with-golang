## Requirements

### Requirement: scanner 套件公開 DetectionEngine 介面
`scanner` 套件應公開 `DetectionEngine` 介面，定義可插拔偵測後端的合約。介面應包含兩個方法：`Name() string`（回傳穩定的引擎識別名稱）與 `Inspect(ctx context.Context, path string) ([]EngineDetection, error)`（對指定檔案執行偵測並回傳零或多筆結果）。

#### Scenario: 實作介面的引擎可被 Scan() 使用
- **WHEN** 任意型別實作了 `DetectionEngine` 介面
- **THEN** 可放入 `Options.ExtraEngines` 並由 `Scan()` 在每個符合條件的檔案上呼叫其 `Inspect` 方法

#### Scenario: Inspect 回傳非 nil error 時計入 ErrorFiles
- **WHEN** `Inspect` 對某個檔案回傳非 nil 的 error
- **THEN** 該檔案計入 `ScanReport.ErrorFiles`，掃描繼續處理其他檔案

### Requirement: scanner 套件公開 EngineDetection 型別
`scanner` 套件應公開 `EngineDetection` 結構，包含三個字串欄位：`Name`（威脅名稱）、`Category`（分類）、`Severity`（嚴重度）。此型別為 `DetectionEngine.Inspect` 的回傳元素，由掃描器轉換為 `Detection` 結構後加入報告。

#### Scenario: EngineDetection 欄位正確對應至 Detection
- **WHEN** `Inspect` 回傳含 `Name`、`Category`、`Severity` 的 `EngineDetection`
- **THEN** 掃描器將其轉換為 `Detection`，其中 `Name`、`Category`、`Severity` 值完整保留，`Engine` 欄位設為該引擎的 `Name()` 回傳值，`SHA256` 欄位設為該檔案已計算的雜湊值

### Requirement: Options 支援 ExtraEngines 與 FileTimeout
`scanner.Options` 應包含 `ExtraEngines []DetectionEngine` 欄位（額外引擎清單；nil 或空清單表示僅使用 hash 引擎）與 `FileTimeout time.Duration` 欄位（每個檔案的 context 截止時間；0 表示不設限）。

#### Scenario: ExtraEngines 為 nil 時行為與原有 hash-only 模式相同
- **WHEN** `Options.ExtraEngines` 為 nil 且 `Options.FileTimeout` 為 0
- **THEN** `Scan()` 的行為與未加入此功能前完全相同

#### Scenario: FileTimeout 傳遞至每個引擎的 Inspect 呼叫
- **WHEN** `Options.FileTimeout` 設為 N 秒
- **THEN** 每個檔案的 `Inspect` 呼叫使用帶有 N 秒截止時間的 `context.Context`

### Requirement: Detection 結構包含 Engine 欄位
`scanner.Detection` 結構應包含 `Engine string` 欄位，標示產生此偵測結果的引擎名稱（例如 `"hash"` 或 `"yara"`）。hash 引擎產生的偵測結果應將 `Engine` 設為 `"hash"`。

#### Scenario: hash 引擎偵測結果帶有 Engine="hash"
- **WHEN** 檔案的 SHA256 命中特徵資料庫
- **THEN** 對應的 `Detection.Engine` 值為 `"hash"`

#### Scenario: 同一檔案可同時產生多個引擎的偵測結果
- **WHEN** 某個檔案的 SHA256 命中 hash 引擎且同時符合 YARA 規則
- **THEN** `ScanReport.Detections` 包含兩筆獨立的 `Detection`，分別帶有 `Engine="hash"` 與 `Engine="yara"`

### Requirement: Scan() 接受 context.Context 第一參數
`scanner.Scan()` 函式簽名應為 `Scan(ctx context.Context, db *sigdb.DB, opts Options) (*ScanReport, error)`。呼叫端應傳入 `context.Background()` 或帶有截止時間的 context。

#### Scenario: 傳入 context.Background() 時行為正常
- **WHEN** 以 `context.Background()` 呼叫 `Scan()`
- **THEN** 掃描正常完成，行為與原有實作相同
