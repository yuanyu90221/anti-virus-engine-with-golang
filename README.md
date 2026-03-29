# avengine — Go 防毒引擎

基於 SHA256 雜湊比對與 YARA 規則引擎的輕量級命令列防毒掃描工具。
使用 YAML 格式的特徵資料庫，可整合至 CI/CD 流程，以 Pants 2.30.0 管理建置。

---

## 目錄

- [架構設計](#架構設計)
- [套件說明](#套件說明)
- [資料流程](#資料流程)
- [快速開始](#快速開始)
- [YARA 引擎](#yara-引擎)
- [特徵 YAML 格式](#特徵-yaml-格式)
- [病毒特徵 hash 參考表](#病毒特徵-hash-參考表)
- [結束碼](#結束碼)
- [參考資源](#參考資源)

---

## 架構設計

```
cmd/avengine/
  main.go              ← CLI 入口：解析旗標，串接各套件

internal/
  config/
    config.go          ← CLIConfig 結構體與預設值常數
  hasher/
    hasher.go          ← 串流 SHA256 計算
  sigdb/
    sigdb.go           ← Loader 介面、DB 記憶體索引、Lookup
    loader_yaml.go     ← YAMLLoader：從目錄讀取 .yaml 特徵檔
  scanner/
    scanner.go         ← DetectionEngine 介面、遞迴掃描、多引擎協調
  yara/
    yara.go            ← YARA CLI subprocess 引擎
  reporter/
    reporter.go        ← 文字表格 / JSON 輸出，結束碼常數

signatures/
  ransomware.yaml      ← 勒索軟體特徵（範例）
  trojans.yaml         ← 木馬特徵（範例）

rules/
  ransomware.yar       ← WannaCry、Ryuk、LockBit 特徵規則
  trojan.yar           ← RAT／後門特徵規則（njRAT、AsyncRAT、Gh0st）
  dropper.yar          ← Dropper／Downloader 行為特徵規則
  webshell.yar         ← PHP／ASP／JSP Webshell 特徵規則
  coinminer.yar        ← XMRig 與挖礦池連線特徵規則

testdata/
  clean/harmless.txt        ← 不命中任何特徵的測試檔
  infected/fake_malware.bin ← 雜湊與 ransomware.yaml 相符的測試檔
  yara-e2e/
    malware_sample.bin ← 含 C2 beacon 特徵字串的端對端測試樣本
    clean_file.txt     ← 乾淨的端對端測試檔案
    rules.yar          ← 端對端測試用 YARA 規則
```

### 偵測引擎架構

avengine 支援多個可插拔的偵測後端，透過 `DetectionEngine` 介面統一協調：

```
Scan(ctx, db, opts)
  ├─ [永遠執行] hash 引擎 → SHA256 → db.Lookup() → Detection{Engine:"hash"}
  └─ [--yara-rules 時啟用] YARA 引擎 → yara CLI subprocess → Detection{Engine:"yara"}
```

每個偵測結果帶有 `Engine` 欄位，標示命中來源。同一個檔案可同時被多個引擎偵測，產生獨立的偵測記錄。

### 設計原則

| 原則 | 實作方式 |
|------|----------|
| **依賴注入** | `sigdb.DB` 透過 `Loader` 介面取得資料，與磁碟 / 網路等來源解耦 |
| **可插拔引擎** | `DetectionEngine` 介面讓任意偵測後端可接入 `scanner.Scan()`，無需修改核心邏輯 |
| **記憶體效率** | `hasher` 以 `io.Copy` 串流計算，無論檔案多大僅使用固定緩衝區 |
| **錯誤隔離** | `scanner` 遇到單一檔案錯誤只計入 `ErrorFiles`，不中止整體掃描 |
| **輸出解耦** | `reporter.Reporter` 介面讓 main.go 不需知道 text/JSON 的實作細節 |
| **向下相容** | 不提供 `--yara-rules` 時行為與原有 hash-only 模式完全相同 |
| **結束碼語意** | 0=乾淨 / 1=威脅 / 2=錯誤，讓 CI/CD 腳本可直接以 `$?` 判斷 |

---

## 套件說明

### `internal/hasher`

提供 `HashFile(path string) (string, error)`，回傳 64 字元小寫十六進位 SHA256。
使用 `io.Copy` 將檔案資料分塊（預設 32 KB）寫入 `sha256.New()`，不將整個檔案讀入記憶體。

### `internal/sigdb`

**`Loader` 介面**：任何實作 `Load() ([]Signature, error)` 的型別都可作為資料來源。

**`DB`**：啟動時將所有特徵建成 `map[sha256hex]MatchResult` 記憶體索引，
`Lookup(hash)` 為 O(1) 查詢，一旦建立後唯讀，並行讀取安全。

**`YAMLLoader`**：讀取指定目錄中所有 `.yaml` 檔案，合併後回傳；
檔案層級的 `category` 欄位會自動注入每筆 `Signature.Category`。

### `internal/scanner`

以 `filepath.WalkDir` 深度優先遍歷目標目錄。
每個非目錄項目依序：

1. 符號連結檢查：`FollowLinks=false`（預設）時略過，防止迴圈
2. 大小過濾：`MaxFileSizeB > 0` 時略過超大檔案，避免長時間阻塞
3. `hasher.HashFile` 計算 SHA256
4. `db.Lookup` 查詢特徵索引，命中則加入 `Detections`（`Engine: "hash"`）
5. 依序呼叫 `ExtraEngines[i].Inspect(ctx, path)`，將結果加入 `Detections`
6. 若 `Options.OnProgress` 不為 nil，每處理完一個檔案後回呼通知

**`DetectionEngine` 介面**：

```go
type DetectionEngine interface {
    Name() string
    Inspect(ctx context.Context, path string) ([]EngineDetection, error)
}
```

**`Options` 主要欄位**：

| 欄位 | 說明 |
|------|------|
| `Dir` | 遞迴掃描的根目錄 |
| `FollowLinks` | 是否追蹤符號連結 |
| `MaxFileSizeB` | 超過此大小的檔案略過（0 = 不限制） |
| `OnProgress` | 每檔完成後的進度回呼 |
| `ExtraEngines` | 額外偵測引擎清單（nil = hash-only） |
| `FileTimeout` | 每個檔案的 context 截止時間（0 = 不限制） |

### `internal/yara`

透過 YARA CLI subprocess 執行規則比對，實作 `DetectionEngine` 介面。

**`New(rulesPath string) (*Engine, error)`**：以 `exec.LookPath` 查找 `yara` binary，找不到時回傳錯誤（呼叫端可降級為 hash-only）。

**`NewWithBinary(binaryPath, rulesPath string) *Engine`**：指定 binary 路徑，供測試注入假 binary。

**`Inspect(ctx, filePath)`** 執行邏輯：
- 呼叫 `yara <rulesPath> <filePath>`
- exit 0：解析 stdout，每行 `RuleName /path` 對應一筆偵測
- exit 1：無比對，回傳空切片（非錯誤）
- exit 2+：回傳 error，包含 exit code 與 stderr
- context timeout：回傳 error，子程序被終止

### `internal/reporter`

工廠函式 `New(format)` 回傳 `Reporter` 介面實作：

- **`textReporter`**：中文摘要行 + 固定寬度欄位表格（含引擎欄位），適合終端機閱讀
- **`jsonReporter`**：camelCase 鍵名的縮排 JSON（含 `engine` 欄位），時間欄位使用 RFC 3339 格式

### `internal/config`

僅含 `CLIConfig` 結構體與預設值常數，無業務邏輯。

---

## 資料流程

```
[磁碟] signatures/*.yaml
          │
          ▼
  sigdb.YAMLLoader.Load()
          │  解析 YAML，注入 category
          ▼
  sigdb.NewDB()
          │  建立 map[sha256]MatchResult
          ▼
  ┌────────────────────────────────────────────────────────────────┐
  │  scanner.Scan(ctx, db, opts)                                   │
  │                                                                │
  │  filepath.WalkDir(opts.Dir)                                    │
  │    └─ 每個檔案                                                  │
  │         ├─ 符號連結？→ 略過（FollowLinks=false 時）             │
  │         ├─ 超大？   → 略過（MaxFileSizeB > 0 時）               │
  │         ├─ hasher.HashFile()  → SHA256                         │
  │         ├─ db.Lookup(hash)   → 命中 → Detection{Engine:"hash"} │
  │         └─ ExtraEngines[i].Inspect(ctx, path)                  │
  │              └─ yara.Engine  → Detection{Engine:"yara"}        │
  └────────────────────────────────────────────────────────────────┘
          │
          ▼
  ScanReport { Detections[]{Path,SHA256,Engine,...}, TotalFiles, ... }
          │
          ▼
  reporter.Write(stdout, report)
          │
          ├─ text → 中文表格（含「引擎」欄位）
          └─ json → camelCase JSON（含 "engine" 鍵）

  os.Exit(0 | 1 | 2)

[進度顯示（text 模式 + TTY）]

  scanner.Scan(ctx, db, opts)
    └─ OnProgress(path, count)      ← 每檔回呼
          │
          ▼
  stderr: \r[N] path                ← 同行覆寫，不影響 stdout
          │
          ▼（掃描完成）
  stderr: \r\033[K                  ← 清除進度列
          │
          ▼
  stdout: 報告輸出
```

---

## 快速開始

### 環境需求

- Go 1.22+（本專案使用 Go 1.25.0）
- Pants 2.30.0（使用系統已安裝的 `pants`，或依下方步驟下載）
- YARA 4.0+（選用，僅啟用 `--yara-rules` 時需要，安裝方式見 [YARA 引擎](#yara-引擎)）

### 建置

```bash
# （選擇性）下載 Pants scie-pants 啟動器至本機
curl -fsSL https://pantsbuild.github.io/setup/pants -o ./pants && chmod +x ./pants

# 安裝相依套件
go mod download

# 自動產生 BUILD 檔案（首次或新增套件後執行）
pants tailor ::

# 執行所有測試
pants test ::

# 建置二進位檔（輸出至 dist/cmd.avengine/bin）
pants package cmd/avengine:
```

### 執行

```bash
# 掃描目錄（僅 hash 引擎，終端機下自動顯示即時進度）
./dist/cmd.avengine/bin scan --dir ./testdata --sigs ./signatures

# JSON 輸出（適合 CI/CD 整合，不顯示進度）
./dist/cmd.avengine/bin scan --dir ./testdata --sigs ./signatures --output json

# 同時啟用 YARA 引擎
./dist/cmd.avengine/bin scan --dir ./testdata --sigs ./signatures \
  --yara-rules ./rules.yar

# 略過超過 10 MB 的檔案，並追蹤符號連結
./dist/cmd.avengine/bin scan --dir /path/to/scan --sigs ./signatures \
  --max-size 10 --follow-links
```

**即時進度顯示**：在互動式終端機（TTY）以 `text` 模式執行時，掃描過程會於 stderr 顯示單行進度，格式為：

```
[42] /path/to/current/file.bin
```

掃描完成後進度列自動清除，不影響 stdout 的報告輸出。JSON 模式或非 TTY 環境（如 CI pipeline）下不顯示進度。

### 所有旗標

| 旗標 | 預設值 | 說明 |
|------|--------|------|
| `--dir` | （必填）| 遞迴掃描的目標目錄 |
| `--sigs` | `./signatures` | 特徵 YAML 目錄 |
| `--output` | `text` | 輸出格式：`text` 或 `json` |
| `--follow-links` | `false` | 追蹤符號連結 |
| `--max-size` | `0`（不限制） | 略過超過 N MB 的檔案 |
| `--verbose` | `false` | 顯示所有掃描結果（含乾淨檔案，預留功能） |
| `--yara-rules` | `""`（停用） | YARA 規則檔案或目錄路徑 |
| `--yara-timeout` | `10` | 每個檔案的 YARA subprocess 逾時（秒） |

---

## YARA 引擎

YARA 是業界標準的惡意軟體規則語言，透過字串、正規表達式、位元組模式及 PE 結構分析偵測威脅。hash 引擎只能比對已知精確樣本，YARA 則可涵蓋變種與未入庫的威脅。

### 安裝 YARA

**Debian / Ubuntu**

```bash
sudo apt-get install -y yara
yara --version  # 應輸出 4.x.x
```

**macOS（Homebrew）**

```bash
brew install yara
yara --version
```

**從原始碼編譯（Alpine / 自訂 CI 映像）**

```bash
sudo apt-get install -y automake libtool make gcc pkg-config \
    libssl-dev libjansson-dev libmagic-dev

curl -L https://github.com/VirusTotal/yara/archive/refs/tags/v4.5.2.tar.gz | tar xz
cd yara-4.5.2
./bootstrap.sh
./configure --with-crypto --enable-magic --enable-dotnet
make && sudo make install
```

### 驗證 YARA 安裝

```bash
# 建立測試規則
cat > /tmp/test.yar << 'EOF'
rule HelloWorld {
    strings:
        $a = "Hello, YARA"
    condition:
        $a
}
EOF

# 建立測試檔案
echo "Hello, YARA" > /tmp/sample.txt

# 驗證比對（應輸出 "HelloWorld /tmp/sample.txt"）
yara /tmp/test.yar /tmp/sample.txt

# 驗證無比對（應無輸出，結束碼 1）
echo "no match here" > /tmp/clean.txt
yara /tmp/test.yar /tmp/clean.txt; echo "exit: $?"
```

### 撰寫 YARA 規則

```yara
rule SuspiciousC2 {
    meta:
        description = "Detects C2 beacon patterns"
        severity     = "high"
    strings:
        $c2  = "C2_SERVER="
        $key = "XOR_KEY=0x"
    condition:
        $c2 and $key
}
```

規則檔案可為單一 `.yar` 檔案，或包含多個 `.yar` 檔案的目錄。

### 內建規則集（`rules/`）

專案內建的 `rules/` 目錄提供涵蓋常見威脅的開箱即用規則，依類別分檔存放：

| 檔案 | 涵蓋威脅 | 嚴重度 |
|------|----------|--------|
| `ransomware.yar` | WannaCry、Ryuk、LockBit | high |
| `trojan.yar` | njRAT、AsyncRAT、Gh0st RAT 及通用 RAT 特徵 | high |
| `dropper.yar` | PowerShell 下載器、Mshta dropper、通用解碼執行 | medium–high |
| `webshell.yar` | PHP eval+base64、China Chopper、ASP CreateObject、JSP Runtime.exec | high |
| `coinminer.yar` | XMRig、常見挖礦池 URL、瀏覽器端 CryptoNight 腳本 | medium |

每條規則均包含 `description`、`severity`、`date`、`reference` 等 meta 欄位，方便稽核與更新。

```bash
# 載入全部內建規則掃描目標目錄
avengine scan --dir ./targets --sigs ./signatures --yara-rules ./rules/
```

### 使用範例

```bash
# 以 YARA 規則掃描（text 輸出）
avengine scan --dir ./targets --sigs ./signatures --yara-rules ./rules/

# JSON 輸出，確認 engine 欄位
avengine scan --dir ./targets --sigs ./signatures \
  --yara-rules ./rules.yar --output json | jq '.detections[].engine'

# 設定每個檔案最多 30 秒逾時
avengine scan --dir ./targets --sigs ./signatures \
  --yara-rules ./rules.yar --yara-timeout 30
```

**輸出範例（text 模式）**

```
掃描完成
檔案總數: 3  威脅: 2  錯誤: 0  略過: 0  耗時: 45.2ms

路徑                          SHA256(前16)      引擎    威脅名稱          嚴重度   分類
------------------------------------------------------------------------
targets/sample.bin            74c7308f2d7debda  yara    SuspiciousC2  unknown  yara
targets/known.bin             91a3f84e7eef3bd8  hash    WannaCry      critical ransomware
```

**yara binary 找不到時的降級行為**

```
warning: YARA engine unavailable: yara: binary not found in PATH: ...
掃描完成（hash-only 模式）
```

---

## 特徵 YAML 格式

特徵資料庫為 YAML 檔案，放置於 `--sigs` 指定的目錄下。引擎啟動時讀取目錄中所有 `.yaml` 並合併為單一索引。

```yaml
version: "1.0"
category: "ransomware"   # 分類名稱，會注入到每筆特徵的 category 欄位
updated: "2024-01-15"

signatures:
  - sha256: "275a021bbfb6489e54d471899f7db9d1663fc695ec2fe2a2c4538aabf651fd0f"
    name: "EICAR Test File"
    severity: "low"       # low | medium | high | critical
    added: "2024-01-01"
  - sha256: "24d004a104d4d54034dbcffc2a4b19a11f39008a575aa614ea04703480b1022c"
    name: "WannaCry Ransomware"
    severity: "critical"
    added: "2024-01-01"
```

**欄位說明**

| 欄位 | 層級 | 說明 |
|------|------|------|
| `version` | 檔案 | 特徵庫版本號（目前未強制驗證） |
| `category` | 檔案 | 惡意軟體分類，自動注入每筆特徵 |
| `updated` | 檔案 | 最後更新日期（YYYY-MM-DD） |
| `sha256` | 特徵 | 64 字元小寫十六進位，為主要比對鍵 |
| `name` | 特徵 | 威脅名稱（顯示於報告） |
| `severity` | 特徵 | 嚴重程度 |
| `added` | 特徵 | 加入日期 |

---

## 病毒特徵 hash 參考表

下列 SHA256 均來自公開已記錄的惡意軟體研究資料，僅供教育與測試用途。
這些 hash 值為樣本識別碼，不含任何可執行程式碼。

| 名稱 | SHA256 | 分類 | 嚴重程度 | 來源 |
|------|--------|------|----------|------|
| EICAR Test File | `275a021bbfb6489e54d471899f7db9d1663fc695ec2fe2a2c4538aabf651fd0f` | 測試 | low | [eicar.org](https://www.eicar.org/) |
| WannaCry | `24d004a104d4d54034dbcffc2a4b19a11f39008a575aa614ea04703480b1022c` | 勒索軟體 | critical | [CISA Alert AA17-132A](https://www.cisa.gov/news-events/cybersecurity-advisories/aa17-132a) |
| NotPetya | `027cc450ef5f8c5f653329641ec1fed91f694e0d229928963b30f6b0d7d3a745` | 勒索軟體 | critical | [US-CERT Alert TA17-181A](https://www.cisa.gov/news-events/alerts/2017/06/30/petya-ransomware) |
| Mirai | `9a024b9ef95a1ed9e5acba3e2fe2427395e866c42b5bce04d35ca2cefd8d2e4d` | 殭屍網路 | high | [Malware Traffic Analysis](https://www.malware-traffic-analysis.net/) |

---

## 結束碼

| 代碼 | 常數 | 意義 |
|------|------|------|
| `0` | `ExitClean` | 掃描完成，未發現威脅 |
| `1` | `ExitDetected` | 偵測到至少一個威脅（hash 或 YARA） |
| `2` | `ExitError` | 執行錯誤（特徵目錄不存在、缺少必要參數等） |

**CI/CD 整合範例**

```bash
pants package cmd/avengine:

./dist/cmd.avengine/bin scan \
  --dir ./dist \
  --sigs ./signatures \
  --yara-rules ./rules/ \
  --output json | tee scan-report.json

if [ $? -eq 1 ]; then
  echo "::error::掃描偵測到威脅，阻止部署"
  exit 1
fi
```

---

## 參考資源

- YouTube 教學影片：[Building an Antivirus Engine](https://www.youtube.com/watch?v=s_M1vKp69hA)
- [EICAR 測試檔案標準](https://www.eicar.org/download-anti-malware-testfile/)
- [YARA 官方文件](https://yara.readthedocs.io/)
- [YARA GitHub](https://github.com/VirusTotal/yara)
- [Pants Build 官方文件](https://www.pantsbuild.org/docs)
- [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3)
