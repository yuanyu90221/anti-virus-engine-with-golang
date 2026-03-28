### Requirement: pants.toml 設定 Go 後端
專案根目錄必須存在 `pants.toml`，啟用 `pants.backend.experimental.go` 並設定版本為 2.23.0，將原始碼根目錄設為 `/`，並設定 `go_search_paths` 以自動探索系統上的 Go 安裝位置。

#### Scenario: pants 版本檢查通過
- **WHEN** 使用者在儲存庫根目錄執行 `./pants --version`
- **THEN** 輸出內容必須為 `2.23.0`

#### Scenario: pants 列出 Go 目標
- **WHEN** 使用者執行 `./pants help goals`
- **THEN** 輸出內容必須包含 `test`、`package` 與 `tailor` 目標

### Requirement: go.mod 初始化模組
專案根目錄必須存在 `go.mod` 檔案，宣告模組為 `github.com/yuanyu/avengine`，並要求 Go 1.22 以上版本。

#### Scenario: go.mod 存在且有效
- **WHEN** 使用者執行 `go mod verify`
- **THEN** 該指令必須以代碼 0 結束

### Requirement: 根目錄 BUILD 檔案存在
專案根目錄必須存在 `BUILD` 檔案，以便 Pants 能夠探索專案結構。

#### Scenario: tailor 執行無誤
- **WHEN** 至少存在一個 `.go` 原始碼檔案，且使用者執行 `./pants tailor ::`
- **THEN** 該指令必須以代碼 0 結束，並為 Go 套件產生 `BUILD` 條目
