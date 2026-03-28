## ADDED Requirements

### Requirement: Loader 介面允許抽換特徵來源格式
系統 SHALL 定義 `Loader` 介面（`Load() ([]Signature, error)`），使任何格式的特徵載入器均可插入 `DB`，不需修改核心索引邏輯。

#### Scenario: 注入自訂 Loader
- **WHEN** 呼叫 `NewDB(loader)` 並傳入實作 `Loader` 介面的物件
- **THEN** DB 以該 Loader 回傳的特徵建立記憶體索引，且不依賴具體格式

### Requirement: YAMLLoader 從目錄載入所有 .yaml 特徵檔
系統 SHALL 提供 `YAMLLoader{Dir: "<path>"}` 實作，讀取指定目錄下所有 `.yaml` 檔案，解析其 `signatures` 陣列，並回傳 `[]Signature`。YAML 格式如下：`version`、`category`、`updated`、`signatures[]`（含 `sha256`、`name`、`severity`、`added`）。

#### Scenario: 目錄含多個 YAML 檔案
- **WHEN** 目錄下存在 `ransomware.yaml` 與 `trojans.yaml`
- **THEN** `Load()` 回傳兩個檔案所有特徵的合併 slice，且 error 為 nil

#### Scenario: 目錄不存在
- **WHEN** 指定的目錄路徑不存在
- **THEN** `Load()` 回傳 nil 與非 nil error，錯誤訊息包含路徑

#### Scenario: YAML 格式錯誤
- **WHEN** 目錄中某個 `.yaml` 檔案格式不合法
- **THEN** `Load()` 回傳 nil 與非 nil error，錯誤訊息包含問題檔案名稱

### Requirement: DB 以 SHA256 為鍵建立記憶體索引並提供查詢
系統 SHALL 在 `NewDB()` 時將所有 `Signature` 建立為 `map[string]MatchResult`（鍵為小寫 SHA256 hex），並提供 `Lookup(sha256hex string) (MatchResult, bool)` 方法。

#### Scenario: 查詢已知惡意 hash
- **WHEN** 呼叫 `db.Lookup(hash)` 且 hash 存在於索引中
- **THEN** 回傳對應的 `MatchResult`（含 Name、Category、Severity），以及 `true`

#### Scenario: 查詢未知 hash
- **WHEN** 呼叫 `db.Lookup(hash)` 且 hash 不在索引中
- **THEN** 回傳零值 `MatchResult` 與 `false`
