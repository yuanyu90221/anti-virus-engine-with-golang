## MODIFIED Requirements

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
