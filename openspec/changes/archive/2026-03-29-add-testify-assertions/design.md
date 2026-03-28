## Context

專案目前有五個測試套件，全部使用 Go 標準庫的 `testing.T` 手寫斷言（`t.Fatal`、`t.Errorf`、`if err != nil`）。
這種方式需要大量樣板程式碼，且失敗時只顯示開發者自行撰寫的訊息，缺乏自動的 expected/actual diff。

`testify` 是 Go 生態系中最廣泛採用的測試斷言函式庫，提供：
- `require`：斷言失敗時立即停止該測試（相當於原 `t.Fatal`）
- `assert`：斷言失敗後繼續執行（相當於原 `t.Error`），收集所有失敗
- 自動產生 `expected: X, actual: Y` 格式的差異訊息

## Goals / Non-Goals

**Goals:**
- 將所有測試套件的斷言統一改用 `testify/require` 與 `testify/assert`
- 前置條件（setup 錯誤）使用 `require`，比對驗證使用 `assert`
- 測試覆蓋範圍與邏輯維持不變

**Non-Goals:**
- 新增測試案例或改變現有測試邏輯
- 使用 `testify/mock`（本專案以介面 + 假實作取代 mock）
- 使用 `testify/suite`（現有測試規模不需要 suite 組織）

## Decisions

**決策 1：`require` vs `assert` 的使用規範**

- `require`：用於測試的前置條件，如 `require.NoError(t, err)` 在 setup 失敗時立即停止，避免後續 nil dereference
- `assert`：用於行為驗證，如 `assert.Equal(t, expected, actual)`，允許一次執行收集多個失敗

替代方案考慮：僅使用 `require`（更簡單但一個失敗就停止，遮蔽其他問題）→ 採用雙軌制更靈活。

**決策 2：import alias**

使用標準 alias：
```go
import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)
```
不使用 dot import（`import . "..."`)，避免污染命名空間。

**決策 3：Pants BUILD 檔案**

Pants 的 `go_test` 目標透過 `pants tailor ::` 自動從 import 追蹤相依，
加入 testify import 後執行 `pants tailor ::` 即可，無需手動修改 BUILD 檔案。

## Risks / Trade-offs

- [風險] `go mod tidy` 後 `go.sum` 會增加 testify 及其間接相依（`github.com/davecgh/go-spew`、`github.com/pmezard/go-difflib`、`gopkg.in/yaml.v3` 已存在）→ **緩解**：這些都是純開發相依，不進入生產二進位
- [風險] `pants tailor ::` 可能重新產生 BUILD 檔案並覆蓋手動調整 → **緩解**：本專案 BUILD 檔案均由 tailor 管理，無手動修改

## Migration Plan

1. `go get github.com/stretchr/testify` 加入相依
2. `go mod tidy` 清理並更新 `go.sum`
3. 逐一修改五個 `_test.go` 檔案
4. `pants test ::` 驗證全部通過
5. 執行 `pants tailor ::` 確認 BUILD 檔案同步（如有需要）
