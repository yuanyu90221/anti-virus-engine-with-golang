## 1. 專案骨架

- [x] 1.1 安裝 Pants：`curl -fsSL https://static.pantsbuild.org/setup/pants > pants && chmod +x pants`
- [x] 1.2 建立 `pants.toml`（pants_version = "2.30.0"，backend_packages = ["pants.backend.experimental.go"]，go_search_paths 指向 Go 1.25.0 安裝路徑）
- [x] 1.3 執行 `go mod init github.com/yuanyu90221/avengine` 產生 `go.mod`
- [x] 1.4 執行 `go get gopkg.in/yaml.v3` 加入唯一外部相依並產生 `go.sum`
- [x] 1.5 建立根目錄 `BUILD` 檔（可為空或含 `go_mod_tidy` 目標）

## 2. internal/config

- [x] 2.1 建立 `internal/config/config.go`，定義 `CLIConfig` 結構體（Dir、SigsDir、Output、FollowLinks、MaxSizeMB、Verbose）與預設值常數

## 3. internal/hasher

- [x] 3.1 建立 `internal/hasher/hasher.go`，實作 `HashFile(path string) (string, error)`，以 `io.Copy` 串流讀取至 `sha256.New()`
- [x] 3.2 建立 `internal/hasher/hasher_test.go`，涵蓋：正常檔案、空檔案、不存在路徑三個情境

## 4. internal/sigdb — 核心介面與 DB

- [x] 4.1 建立 `internal/sigdb/sigdb.go`，定義 `Loader` 介面、`Signature`、`MatchResult`、`DB` 結構體
- [x] 4.2 在 `sigdb.go` 中實作 `NewDB(loader Loader) (*DB, error)` 與 `(db *DB) Lookup(sha256hex string) (MatchResult, bool)`
- [x] 4.3 建立 `internal/sigdb/sigdb_test.go`，涵蓋：Lookup 命中、Lookup 未命中兩個情境

## 5. internal/sigdb — YAMLLoader

- [x] 5.1 建立 `internal/sigdb/loader_yaml.go`，定義 `YAMLLoader{Dir string}` 並實作 `Load() ([]Signature, error)`
- [x] 5.2 `Load()` 應讀取 `Dir` 下所有 `.yaml` 檔，解析 `version`/`category`/`signatures` 欄位，並合併回傳
- [x] 5.3 建立 `internal/sigdb/loader_yaml_test.go`，涵蓋：多檔案合併、目錄不存在、YAML 格式錯誤三個情境

## 6. internal/scanner

- [x] 6.1 建立 `internal/scanner/scanner.go`，定義 `Options`（Dir、FollowLinks、MaxFileSizeB）與 `FileResult`（Path、SHA256、Matched bool、MatchResult、Skipped bool、Err error）
- [x] 6.2 在 `scanner.go` 中定義 `ScanReport`（Detections、TotalFiles、ErrorFiles、SkippedFiles、StartedAt、FinishedAt）
- [x] 6.3 實作 `Scan(db *sigdb.DB, opts Options) (*ScanReport, error)`：使用 `fs.WalkDir`，處理符號連結略過、大小限制略過、hash 計算與 DB 查詢
- [x] 6.4 建立 `internal/scanner/scanner_test.go`，涵蓋：乾淨目錄、含惡意檔案、符號連結略過、超大檔案略過、無讀取權限五個情境

## 7. internal/reporter

- [x] 7.1 建立 `internal/reporter/reporter.go`，定義常數 `ExitClean = 0`、`ExitDetected = 1`、`ExitError = 2`
- [x] 7.2 定義 `Reporter` 介面（`Write(w io.Writer, report *scanner.ScanReport) error`）與 `New(format string) Reporter`
- [x] 7.3 實作 `textReporter`：輸出掃描摘要表頭與偵測結果表格（路徑、SHA256 前 16 字元、名稱、嚴重程度、分類）
- [x] 7.4 實作 `jsonReporter`：將 `ScanReport` 以 `encoding/json` 序列化為 camelCase JSON
- [x] 7.5 建立 `internal/reporter/reporter_test.go`，涵蓋：text 無威脅、text 有威脅、json 有威脅、無效 format 四個情境

## 8. cmd/avengine/main.go

- [x] 8.1 建立 `cmd/avengine/main.go`，解析 `scan` 子命令與所有旗標（--dir、--sigs、--output、--follow-links、--max-size、--verbose）
- [x] 8.2 缺少 `--dir` 時輸出使用說明至 stderr 並以結束碼 2 退出
- [x] 8.3 串接完整流程：`sigdb.NewDB(YAMLLoader)` → `scanner.Scan` → `reporter.Write` → `os.Exit`
- [x] 8.4 `--max-size` 旗標單位為 MB，轉換為 bytes 後傳入 `scanner.Options.MaxFileSizeB`

## 9. 特徵資料庫範例

- [x] 9.1 建立 `signatures/ransomware.yaml`，包含至少 2 筆勒索軟體特徵（其中一筆 SHA256 需與 testdata 假惡意檔案相符）
- [x] 9.2 建立 `signatures/trojans.yaml`，包含至少 2 筆木馬特徵
- [x] 9.3 建立 `signatures/BUILD`，使用 `resources()` 目標包含所有 `.yaml` 檔案

## 10. 測試資料

- [x] 10.1 建立 `testdata/clean/harmless.txt`（內容隨意，但其 SHA256 不可命中特徵資料庫）
- [x] 10.2 建立 `testdata/infected/fake_malware.bin`（計算其 SHA256 並確保與 `ransomware.yaml` 中某筆特徵相符）

## 11. Pants BUILD 檔與驗證

- [x] 11.1 執行 `./pants tailor ::` 自動產生所有套件的 `BUILD` 檔
- [x] 11.2 執行 `./pants test ::` 確認所有測試通過
- [x] 11.3 執行 `./pants package cmd/avengine:` 產生二進位檔
- [x] 11.4 煙霧測試：`./dist/cmd.avengine/avengine scan --dir ./testdata/infected --sigs ./signatures`（預期結束碼 1）
- [x] 11.5 煙霧測試：`./dist/cmd.avengine/avengine scan --dir ./testdata/clean --sigs ./signatures`（預期結束碼 0）
