### Requirement: CI workflow 在 push 與 PR 時自動觸發
`.github/workflows/ci.yml` 中 SHALL 設定在推送至 `main` 分支及所有 pull request 時自動觸發 CI pipeline。

#### Scenario: push 至 main 時觸發
- **WHEN** 開發者推送 commit 至 `main` 分支
- **THEN** GitHub Actions 自動啟動 CI workflow

#### Scenario: PR 開啟或更新時觸發
- **WHEN** 開發者開啟或更新 pull request
- **THEN** GitHub Actions 自動啟動 CI workflow

### Requirement: CI 執行 pants test 並回報結果
CI workflow SHALL 執行 `pants test ::` 並在測試失敗時以非零結束碼中止 pipeline。

#### Scenario: 所有測試通過
- **WHEN** `pants test ::` 執行完畢且所有測試套件通過
- **THEN** CI job 標記為成功（綠燈）

#### Scenario: 測試失敗時 CI 中止
- **WHEN** `pants test ::` 中任一測試失敗
- **THEN** CI job 標記為失敗（紅燈），不繼續後續步驟

### Requirement: CI 執行 pants package 驗證二進位可建置
CI workflow SHALL 執行 `pants package cmd/avengine:`，驗證二進位檔可成功產生。

#### Scenario: 二進位建置成功
- **WHEN** `pants package cmd/avengine:` 執行完畢
- **THEN** `dist/cmd.avengine/bin` 存在，CI job 繼續執行

#### Scenario: 建置失敗時 CI 中止
- **WHEN** `pants package cmd/avengine:` 因編譯錯誤失敗
- **THEN** CI job 標記為失敗

### Requirement: CI 快取 Pants 與 Go 建置產物
CI workflow SHALL 快取 Pants PEX bootstrap 與 Go module cache，以縮短後續執行時間。

#### Scenario: 快取命中時加速執行
- **WHEN** `go.sum` 與 `pants_version` 均未變更的情況下再次執行 CI
- **THEN** Pants PEX 與 Go modules 從快取還原，不重新下載

### Requirement: CI 使用正確的 Go 版本
CI workflow SHALL 安裝 Go 1.25.0，與 `go.mod` 中指定的版本一致。

#### Scenario: Go 版本與 go.mod 一致
- **WHEN** CI 在 `ubuntu-latest` 上執行
- **THEN** `go version` 輸出 `go1.25.0`，編譯使用此版本
