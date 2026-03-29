## Why

SHA256 hash 比對只能偵測已知的精確樣本，對於變種惡意軟體或未入庫的威脅無能為力。加入 YARA 規則引擎，讓引擎能透過字串、位元組模式及 PE 結構分析偵測未知威脅，大幅提升偵測覆蓋率。

## What Changes

- 新增 `DetectionEngine` 介面與 `EngineDetection` 型別，作為可插拔偵測後端的擴充點
- 新增 `internal/yara` 套件，透過 YARA CLI subprocess 執行規則比對
- `scanner.Options` 新增 `ExtraEngines`（額外引擎清單）與 `FileTimeout`（每檔逾時）欄位
- `scanner.Scan()` 加入 `context.Context` 第一參數，支援逾時傳遞（**BREAKING**：呼叫端需補上 `context.Background()`）
- `scanner.Detection` 新增 `Engine string` 欄位，標示命中來源（"hash" 或 "yara"）
- CLI 新增 `--yara-rules` 與 `--yara-timeout` 旗標；不提供 `--yara-rules` 時行為完全不變
- Reporter 文字表格新增「引擎」欄位；JSON 輸出新增 `engine` 鍵

## Capabilities

### New Capabilities
- `detection-engine-interface`：可插拔偵測引擎介面（`DetectionEngine`、`EngineDetection`）及多引擎掃描流程
- `yara-engine`：YARA CLI subprocess 引擎，含安裝方式、規則格式、subprocess 行為與逾時處理
- `yara-cli-setup`：YARA CLI 工具的安裝方法（apt / Homebrew / 原始碼編譯）與安裝驗證程序

### Modified Capabilities
- `scan-report`：偵測結果新增 `Engine` 欄位，文字與 JSON 輸出格式均需反映此變更
- `avengine-cli`：新增 `--yara-rules` 與 `--yara-timeout` 旗標，以及 `Scan()` 函式簽名變更（加入 `context.Context`）

## Impact

- `internal/scanner/scanner.go`：新增介面型別、修改 `Detection`、`Options`、`Scan()` 簽名
- `internal/yara/yara.go`（新套件）：YARA subprocess 邏輯
- `internal/reporter/reporter.go`：文字與 JSON 輸出新增引擎欄位
- `internal/config/config.go`：新增 YARA 旗標預設值
- `cmd/avengine/main.go`：新增旗標解析與引擎初始化
- 外部依賴：`yara` CLI binary（選用，不在 go.mod 中）
