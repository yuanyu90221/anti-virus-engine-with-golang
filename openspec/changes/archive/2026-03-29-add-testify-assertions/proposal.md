## Why

目前所有測試套件使用標準庫的 `t.Fatal` / `t.Errorf` 手寫斷言，缺乏可讀性且錯誤訊息不夠描述性。
引入 `github.com/stretchr/testify` 可提供語義化的 `assert` / `require` API，讓測試意圖更清晰、失敗訊息更易於診斷。

## What Changes

- 在 `go.mod` 中加入 `github.com/stretchr/testify` 相依套件
- 將四個測試套件的斷言全部改用 `testify/assert`（非致命）與 `testify/require`（致命，相當於原 `t.Fatal`）：
  - `internal/hasher/hasher_test.go`
  - `internal/sigdb/sigdb_test.go`
  - `internal/sigdb/loader_yaml_test.go`
  - `internal/scanner/scanner_test.go`
  - `internal/reporter/reporter_test.go`
- 移除原本的手寫 `if err != nil { t.Fatal(...) }` 等樣板程式碼
- 測試邏輯與覆蓋範圍不變，僅替換斷言方式

## Capabilities

### New Capabilities

- `testify-test-assertions`：以 testify 作為統一的測試斷言函式庫，定義使用 `require` 處理前置條件失敗、`assert` 處理非致命比較的使用規範

### Modified Capabilities

（無）— 測試行為不變，不涉及任何對外 API 或業務邏輯的規格變更

## Impact

- **相依套件**：新增 `github.com/stretchr/testify`（測試用，不影響生產二進位）
- **go.sum**：需執行 `go mod tidy` 更新
- **BUILD 檔案**：Pants `go_test` 目標會自動追蹤新的 import，無需手動修改
- **測試輸出**：testify 在斷言失敗時提供更豐富的 diff 訊息，不影響通過時的行為
