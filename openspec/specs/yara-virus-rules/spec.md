## Requirements

### Requirement: rules 目錄存在並包含 YARA 規則檔案
專案根目錄 SHALL 包含 `rules/` 目錄，其中存放一或多個 `.yar` 格式的規則檔案，依威脅類別分檔組織。

#### Scenario: rules 目錄結構存在
- **WHEN** 使用者 clone 或取得專案後
- **THEN** `rules/` 目錄存在且包含至少以下檔案：`ransomware.yar`、`trojan.yar`、`dropper.yar`、`webshell.yar`、`coinminer.yar`

### Requirement: 每條規則包含完整 meta 資訊
`rules/` 下的每條 YARA 規則 SHALL 包含 `meta` 區段，內含 `description`、`severity`（high/medium/low）、`date`、`reference` 欄位。

#### Scenario: 規則 meta 欄位完整
- **WHEN** 以文字編輯器或程式解析任一 `.yar` 規則
- **THEN** 每條規則的 `meta` 區段均包含 description、severity、date、reference 四個欄位

### Requirement: 規則通過 YARA 語法驗證
所有 `rules/` 下的規則檔案 SHALL 通過 `yara --compile-rules` 語法檢查，不得有語法錯誤。

#### Scenario: 規則語法正確
- **WHEN** 執行 `yara --compile-rules rules/<category>.yar /dev/null`（或等效驗證）
- **THEN** 命令以 exit 0 結束，無語法錯誤輸出

### Requirement: 規則涵蓋常見威脅類別
`rules/` 目錄 SHALL 包含下列威脅類別的規則，每類至少一條：
- Ransomware（勒索軟體）
- Trojan / RAT（木馬 / 遠端存取工具）
- Dropper / Downloader（植入程式）
- Webshell（網頁後門）
- Coinminer（加密貨幣挖礦程式）

#### Scenario: 各類別規則存在
- **WHEN** 列出 `rules/` 目錄內容
- **THEN** 存在對應上述五個類別的 `.yar` 檔案，每個檔案至少包含一條有效規則
