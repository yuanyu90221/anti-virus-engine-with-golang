## 1. 加入相依套件

- [x] 1.1 執行 `go get github.com/stretchr/testify` 將 testify 加入 `go.mod`
- [x] 1.2 執行 `go mod tidy` 更新 `go.sum`，確認間接相依正確解析

## 2. 重構 internal/hasher

- [x] 2.1 修改 `internal/hasher/hasher_test.go`：以 `require.NoError` 取代 setup 的 `t.Fatal`，以 `assert.Equal` / `assert.Len` 取代值比對的 `t.Errorf`
- [x] 2.2 確認 `require` 用於前置條件（`os.CreateTemp`、`f.WriteString`），`assert` 用於 hash 結果驗證

## 3. 重構 internal/sigdb

- [x] 3.1 修改 `internal/sigdb/sigdb_test.go`：以 `require` / `assert` 取代所有手寫斷言
- [x] 3.2 修改 `internal/sigdb/loader_yaml_test.go`：以 `require` / `assert` 取代所有手寫斷言，含多檔案合併、目錄不存在、YAML 錯誤三個情境

## 4. 重構 internal/scanner

- [x] 4.1 修改 `internal/scanner/scanner_test.go`：以 `require` / `assert` 取代所有手寫斷言，涵蓋乾淨目錄、含惡意檔案、符號連結略過、超大檔案略過、無讀取權限五個情境

## 5. 重構 internal/reporter

- [x] 5.1 修改 `internal/reporter/reporter_test.go`：以 `require` / `assert` 取代所有手寫斷言，涵蓋 text/json 輸出與無效格式四個情境

## 6. 驗證

- [x] 6.1 執行 `pants tailor ::` 確認 BUILD 檔案同步 testify import
- [x] 6.2 執行 `pants test ::` 確認四個測試套件全部通過，測試數量與重構前相同
