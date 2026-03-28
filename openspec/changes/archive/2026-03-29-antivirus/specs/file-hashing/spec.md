## ADDED Requirements

### Requirement: 串流計算單一檔案的 SHA256 雜湊
系統 SHALL 以串流方式讀取檔案內容，使用 `crypto/sha256` 計算雜湊，並回傳小寫十六進位字串。實作為 `internal/hasher` 套件中的 `HashFile(path string) (string, error)`。

#### Scenario: 成功計算一般檔案
- **WHEN** 傳入一個可讀取的檔案路徑
- **THEN** 回傳該檔案的 SHA256 十六進位字串（64 字元小寫），且 error 為 nil

#### Scenario: 計算空檔案
- **WHEN** 傳入一個大小為 0 的空檔案路徑
- **THEN** 回傳 `"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"`（SHA256 of empty），且 error 為 nil

#### Scenario: 檔案不存在或無讀取權限
- **WHEN** 傳入一個不存在或無權限的路徑
- **THEN** 回傳空字串，且 error 不為 nil，錯誤訊息包含路徑資訊
