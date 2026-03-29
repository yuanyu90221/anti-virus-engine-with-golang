## MODIFIED Requirements

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
