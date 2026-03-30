## Requirements

### Requirement: yara 套件公開 Engine 型別
`internal/yara` 套件應公開 `Engine` 型別，實作 `scanner.DetectionEngine` 介面。`Engine` 應以 `rulesPath`（YARA 規則檔案或目錄路徑）與 `yaraPath`（yara binary 的絕對路徑）作為內部狀態，建構後為唯讀，並發使用安全。

#### Scenario: Engine 實作 DetectionEngine 介面
- **WHEN** 將 `*yara.Engine` 放入 `scanner.Options.ExtraEngines`
- **THEN** 可被 `scanner.Scan()` 正常呼叫，不發生型別錯誤

### Requirement: New() 驗證 yara binary 可用性
`yara.New(rulesPath string) (*Engine, error)` 應透過 `exec.LookPath("yara")` 解析 binary 路徑。若找不到 binary，回傳描述性 error。成功時回傳已初始化的 `*Engine`。`rulesPath` 可為單一 `.yar` 檔案路徑或包含 `.yar` 檔案的目錄路徑；若為目錄，引擎 SHALL 批次載入目錄下所有 `.yar` 檔案。

#### Scenario: PATH 中存在 yara binary 時 New() 成功
- **WHEN** `yara` binary 存在於系統 PATH 中且 `rulesPath` 為有效路徑（檔案或目錄）
- **THEN** `New()` 回傳非 nil 的 `*Engine` 且 error 為 nil

#### Scenario: PATH 中找不到 yara binary 時 New() 回傳 error
- **WHEN** 系統 PATH 中不存在 `yara` binary
- **THEN** `New()` 回傳 nil 與包含 "not found" 描述的 error

#### Scenario: rulesPath 為目錄時批次載入所有規則
- **WHEN** `rulesPath` 為包含多個 `.yar` 檔案的目錄（例如 `rules/`）
- **THEN** `Inspect()` 可命中目錄下任一規則檔案中定義的規則

### Requirement: Inspect() 呼叫 yara CLI subprocess 並解析輸出
`Engine.Inspect(ctx context.Context, filePath string)` 應執行 `yara --print-meta <rulesPath> <filePath>`，將 stdout 每行解析為 `EngineDetection`。每行格式為 `RuleName [key="val",...] /path/to/file`，`Name` 設為規則名稱，`Category` 預設為 `"yara"`，`Severity` 從 metadata block 中的 `severity` 欄位解析；若規則未定義 `severity` metadata，`Severity` 預設為 `"unknown"`。

#### Scenario: YARA 規則命中時回傳對應偵測結果
- **WHEN** `yara` subprocess 以 exit 0 結束並在 stdout 輸出一行或多行規則命中
- **THEN** `Inspect()` 回傳與命中行數相等的 `EngineDetection` 切片，每筆 `Name` 為對應的規則名稱

#### Scenario: 規則有 severity metadata 時回傳正確嚴重程度
- **WHEN** `yara --print-meta` 輸出的行包含 `severity="high"` 或其他有效值
- **THEN** 對應的 `EngineDetection.Severity` 為該值（例如 `"high"`、`"medium"`）

#### Scenario: 規則無 severity metadata 時 Severity fallback 為 unknown
- **WHEN** `yara --print-meta` 輸出的 metadata block 中不含 `severity` 欄位
- **THEN** 對應的 `EngineDetection.Severity` 為 `"unknown"`

#### Scenario: YARA 規則未命中時回傳空切片
- **WHEN** `yara` subprocess 以 exit 1 結束（YARA CLI 的無比對語意）
- **THEN** `Inspect()` 回傳 nil 切片且 error 為 nil

#### Scenario: YARA subprocess 以 exit 2 以上結束時回傳 error
- **WHEN** `yara` subprocess 以 exit code 2 或更高值結束（規則語法錯誤、檔案無法讀取等）
- **THEN** `Inspect()` 回傳非 nil 的 error，error 訊息包含 exit code 與 stderr 內容

#### Scenario: context 截止時間到期時 Inspect() 回傳 error
- **WHEN** context 在 subprocess 完成前到期
- **THEN** `Inspect()` 回傳非 nil 的 error，子程序被終止

### Requirement: NewWithBinary() 支援測試時注入假 binary
`yara.NewWithBinary(binaryPath, rulesPath string) *Engine` 應直接使用指定的 binary 路徑建立 `Engine`，跳過 `exec.LookPath` 查詢。此函式專供測試使用，讓測試可注入假的 yara binary。

#### Scenario: 注入假 binary 後 Inspect() 使用該路徑執行
- **WHEN** 以假 binary 路徑呼叫 `NewWithBinary` 並執行 `Inspect()`
- **THEN** 執行的是假 binary 而非系統的 yara binary
