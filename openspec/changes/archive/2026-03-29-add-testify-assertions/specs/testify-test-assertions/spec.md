## ADDED Requirements

### Requirement: 使用 testify 作為測試斷言函式庫
專案的所有測試套件 SHALL 使用 `github.com/stretchr/testify` 作為統一的斷言函式庫，
取代手寫的 `if condition { t.Fatal(...) }` 樣板程式碼。

#### Scenario: 前置條件失敗立即停止
- **WHEN** 測試的 setup 步驟（如開啟檔案、建立暫存目錄）回傳 error
- **THEN** 使用 `require.NoError(t, err)` 或 `require.NoError(t, err, ...)` 使測試立即停止，不繼續執行後續步驟

#### Scenario: 行為驗證收集所有失敗
- **WHEN** 測試驗證回傳值、計數或字串等行為結果
- **THEN** 使用 `assert.Equal`、`assert.Len`、`assert.Contains` 等非致命斷言，使單次執行可收集所有失敗訊息

### Requirement: require 用於前置條件，assert 用於行為比對
測試程式碼 SHALL 依照以下規範區分 `require` 與 `assert` 的使用時機：
- `require`：任何若失敗則後續步驟必然 panic 或無意義的前置條件
- `assert`：任何單純的值比對或狀態驗證，失敗後仍可繼續蒐集其他失敗

#### Scenario: require 用於 error 前置條件
- **WHEN** 測試呼叫可能回傳 error 的函式，且後續步驟依賴其回傳值
- **THEN** 使用 `require.NoError(t, err)` 而非 `if err != nil { t.Fatal(err) }`

#### Scenario: assert 用於值驗證
- **WHEN** 測試比對函式的回傳值是否符合預期
- **THEN** 使用 `assert.Equal(t, expected, actual)` 而非 `if got != want { t.Errorf(...) }`

### Requirement: 測試覆蓋範圍與邏輯不變
引入 testify 後，測試案例的覆蓋範圍、輸入情境與驗證邏輯 SHALL 與重構前完全相同，
僅斷言語法改變，不新增或移除測試案例。

#### Scenario: 現有測試案例保持不變
- **WHEN** 以 testify 重構後執行 `pants test ::`
- **THEN** 所有原有測試案例仍然存在且通過，測試數量不減少
