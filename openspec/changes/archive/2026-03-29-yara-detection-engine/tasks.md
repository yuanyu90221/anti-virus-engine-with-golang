## 1. 環境準備：安裝 YARA CLI

- [x] 1.1 確認目標環境（Linux / macOS）並選擇安裝方式
- [x] 1.2 Linux (Debian/Ubuntu)：執行 `sudo apt-get install -y yara`
- [x] 1.3 macOS：執行 `brew install yara`
- [x] 1.4 無套件管理工具環境：依 `yara-cli-setup` spec 從原始碼編譯安裝
- [x] 1.5 執行 `yara --version` 驗證版本 >= 4.0.0

## 2. 驗證 YARA 安裝可正確比對規則

- [x] 2.1 建立測試規則檔案 `test.yar`（含 HelloWorld 規則）
- [x] 2.2 建立包含 `Hello, YARA` 字串的 `sample.txt`
- [x] 2.3 執行 `yara test.yar sample.txt`，確認輸出包含 `HelloWorld` 且結束碼為 0
- [x] 2.4 對不含目標字串的檔案執行比對，確認結束碼為 1（無比對非錯誤）

## 3. 核心實作：DetectionEngine 介面與掃描器整合

- [x] 3.1 `internal/scanner/scanner.go`：新增 `DetectionEngine` 介面與 `EngineDetection` 型別
- [x] 3.2 `internal/scanner/scanner.go`：`Detection` 結構加入 `Engine string` 欄位
- [x] 3.3 `internal/scanner/scanner.go`：`Options` 加入 `ExtraEngines` 與 `FileTimeout` 欄位
- [x] 3.4 `internal/scanner/scanner.go`：`Scan()` 加入 `context.Context` 第一參數，並呼叫額外引擎

## 4. YARA 引擎實作

- [x] 4.1 建立 `internal/yara/yara.go`：`Engine` 型別、`New()`、`NewWithBinary()`
- [x] 4.2 實作 `Inspect()`：呼叫 subprocess、解析 stdout、處理 exit 0/1/2 及 context timeout
- [x] 4.3 建立 `internal/yara/BUILD`：pants 套件定義
- [x] 4.4 建立 `internal/yara/yara_test.go`：假 binary 測試（match、no match、error exit、timeout）

## 5. Reporter 更新

- [x] 5.1 `internal/reporter/reporter.go`：文字表格 header 與每列加入「引擎」欄位
- [x] 5.2 `internal/reporter/reporter.go`：`jsonDetection` 加入 `engine` 欄位
- [x] 5.3 `internal/reporter/reporter_test.go`：更新 `infectedReport()` 加入 `Engine:"hash"`，新增 YARA 偵測結果測試

## 6. CLI 整合

- [x] 6.1 `internal/config/config.go`：加入 `YARARules`、`YARATimeout`、`DefaultYARATimeout` 常數
- [x] 6.2 `cmd/avengine/main.go`：新增 `--yara-rules` 與 `--yara-timeout` 旗標
- [x] 6.3 `cmd/avengine/main.go`：加入 YARA 引擎初始化邏輯（找不到 binary 時 warn 並降級）
- [x] 6.4 `cmd/avengine/main.go`：`Scan()` 呼叫改傳 `context.Background()`、`ExtraEngines`、`FileTimeout`

## 7. 測試更新與驗證

- [x] 7.1 `internal/scanner/scanner_test.go`：全部 `Scan()` 呼叫補上 `context.Background()`
- [x] 7.2 `internal/scanner/scanner_test.go`：新增 `mockEngine` 及多引擎測試（ExtraEngine 偵測、hash+YARA 同時命中）
- [x] 7.3 執行 `go test ./...` 確認全數通過
- [x] 7.4 執行 `pants test ::` 確認所有 pants 目標通過

## 8. 端對端驗證（需系統已安裝 yara）

- [x] 8.1 撰寫含有特定字串的測試惡意樣本與對應 YARA 規則
- [x] 8.2 執行 `avengine scan --dir ./testdata --yara-rules ./rules.yar` 確認 YARA 偵測結果出現在輸出中
- [x] 8.3 執行 JSON 格式掃描，確認 `detections[].engine` 欄位正確顯示 `"yara"`
- [x] 8.4 移除 `--yara-rules` 旗標後重新掃描，確認行為與 hash-only 模式相同
