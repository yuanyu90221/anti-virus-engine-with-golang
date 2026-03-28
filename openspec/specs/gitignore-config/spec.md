## Requirements

### Requirement: Pants 建置產物不被版本控制追蹤
`.gitignore` 中 SHALL 包含 `.pants.d/` 與 `dist/` 的忽略規則，
確保 Pants 快取目錄與建置輸出不被意外提交。

#### Scenario: Pants 快取目錄被忽略
- **WHEN** 執行 `pants test ::` 或 `pants package` 後執行 `git status`
- **THEN** `.pants.d/` 目錄不出現在未追蹤檔案列表中

#### Scenario: 建置輸出目錄被忽略
- **WHEN** 執行 `pants package cmd/avengine:` 後執行 `git status`
- **THEN** `dist/` 目錄不出現在未追蹤檔案列表中

### Requirement: 執行時產生的掃描報告不被版本控制追蹤
`.gitignore` 中 SHALL 包含 `scan-report.json` 的忽略規則，
防止 CI/CD 腳本執行後產生的報告檔被提交。

#### Scenario: 掃描報告被忽略
- **WHEN** 執行 `avengine scan ... | tee scan-report.json` 後執行 `git status`
- **THEN** `scan-report.json` 不出現在未追蹤檔案列表中

### Requirement: IDE 設定目錄不被版本控制追蹤
`.gitignore` 中 SHALL 包含 `.idea/`（JetBrains）與 `.vscode/`（VS Code）的忽略規則，
確保開發者的本機 IDE 設定不汙染共用版本控制。

#### Scenario: IDE 設定目錄被忽略
- **WHEN** 開發者使用 JetBrains 或 VS Code 開啟專案後執行 `git status`
- **THEN** `.idea/` 與 `.vscode/` 不出現在未追蹤檔案列表中

### Requirement: 作業系統暫存檔不被版本控制追蹤
`.gitignore` 中 SHALL 包含 `.DS_Store`（macOS）與 `Thumbs.db`（Windows）的忽略規則。

#### Scenario: macOS 暫存檔被忽略
- **WHEN** 在 macOS 上使用 Finder 瀏覽專案目錄後執行 `git status`
- **THEN** `.DS_Store` 不出現在未追蹤檔案列表中
